package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	. "github.com/woodpecker/tracker"
	"os"
	"strings"
	"strconv"
)

func main() {

	fi, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fi.Close()
	filter := ""
	if len(os.Args) > 2 {
		filter = os.Args[2]
	}

	Readfile(fi, filter)

}

func ReadBytes(fi *os.File, b []byte, needread int) {
	n, err := fi.Read(b)
	if n != needread {
		panic(n)

	}
	if err != nil {
		panic(err)
	}

}

func ReadRecords(fi *os.File, bufflen uint32, filter string) []*Tracker {
	s := strings.Split(filter, "-")
	result := make([]*Tracker, 0)

	timebuff := make([]byte, 4)
	serviceid := make([]byte, 4)
	messagetype := make([]byte, 4)
	datelen := make([]byte, 4)
	var readed uint32
	for readed != bufflen {

		ReadBytes(fi, timebuff, 4)
		ReadBytes(fi, serviceid, 4)
		ReadBytes(fi, messagetype, 4)
		ReadBytes(fi, datelen, 4)

		itimestamp := binary.BigEndian.Uint32(timebuff)

		if filter != "" && (int(itimestamp) != GetIntValue(s[0])) {
			return result
		}

		needread := binary.BigEndian.Uint32(datelen)

		buff := make([]byte, int(needread))

		ReadBytes(fi, buff, int(needread))
		readed += 16 + needread

		tr := &Tracker{}
		err := proto.Unmarshal(buff, tr)
		if err != nil {
			fmt.Println(err)
			return result
		} else {
			if tr.Sid == filter {
				result = append(result, tr)
			}
			if filter == "" {
				fmt.Println(tr)
			}
		}
	}
	return result
}

func Readfile(fi *os.File, filter string) {
	timebuff := make([]byte, 4)
	lenbuff := make([]byte, 4)

	filelen, err := fi.Seek(0, 2)
	fmt.Println("file length:", filelen)
	if err != nil {
		panic(err)
	}
	index := 0
	_, err = fi.Seek(0, 0)

	record := make([]*Tracker, 0)
	for {
		l, err := fi.Seek(0, 1)
		if l == filelen {
			break
		}
		if err != nil {
			panic(err)
		}

		ReadBytes(fi, timebuff, 4)
		ReadBytes(fi, lenbuff, 4)

		needread := binary.BigEndian.Uint32(lenbuff)

		index += 8 + int(needread)
		result := ReadRecords(fi, needread, filter)
		if len(result) > 0 {
			for i := range result {
				fmt.Println(result[i])
				record = append(record, result[i])
			}
		}
		fi.Seek(int64(index), 0)
	}

	if len(record) > 0 {
		CreateJs(record)
	}
}

func CreateJs(result []*Tracker) {
	if len(result) == 0 {
		return
	}
	var mintime, maxtime int64
	var maxdeep int
	for i := range result {
		fmt.Println(result[i].Timestamp)
		if mintime == 0 {
			mintime = result[i].Timestamp
		}
		if maxtime == 0 {
			maxtime = result[i].Timestamp
		}

		if result[i].Timestamp >= maxtime {
			maxtime = result[i].Timestamp
		}

		if result[i].Timestamp <= mintime {
			mintime = result[i].Timestamp
		}

		deepstr := strings.Split(result[i].Callstack, ".")
		if len(deepstr) >= maxdeep {
			maxdeep = len(deepstr)
		}
	}

	timestamp := maxtime - mintime
	unit := timestamp / (26 * 12)

	// need sort the request
	// skip in this version

	result2 := make([]*Tracker, 0)
	for deepth := 0; deepth < maxdeep; deepth++ {
		for {
			result3 := Filter(result, deepth)
			if len(result3) == 0 {
				break
			}
			for i := range result3 {
				if i == len(result3)-1 {
					break
				}
				result3[i].Startime = (result3[i].Timestamp - mintime) / unit
				result3[i].Endtime = (result3[i+1].Timestamp - mintime) / unit
				result2 = append(result2, result3[i])
			}

		}

	}

	fontlist := [6]string{"lorem", "ipsum", "dolor", "ipsum", "default", "sit"}
	htmlstring := make([]string, 0)
	
	for i := range result2 {

		year :=  int( result2[i].Startime/12 + 1900)
		month := int( result2[i].Startime%12 + 1)

		eyear := int( result2[i].Endtime/12 + 1900)
		emonth :=int(  result2[i].Endtime%12 + 1)

		callstack := ""
		if len(result2[i].Callstack) == 0 {
			callstack = result2[i].Callstack + strconv.Itoa(int(result2[i].Id))
		} else {
			callstack = result2[i].Callstack + "." + strconv.Itoa(int(result2[i].Id))
		}

		htmlstring = append( htmlstring, "['" + strconv.Itoa(month) + "/" + strconv.Itoa(year) +  "', '"  + strconv.Itoa(emonth)  + "/" + strconv.Itoa(eyear) +  "', '" + callstack + "-" + result2[i].Function + "', '" +  fontlist[i % len(fontlist)]  + "']" )

	}
	
	printstr := `
		(function(){
	  'use strict';

	  Lib.ready(function() {
	    console.log('ads');

	   // replace code
	   new Timesheet('timesheet-default', 1900, 1927, 
		   [`
	 printstr +=strings.Join(htmlstring, ",\n")
	
	 printstr += `
	 ]);
	 

	    document.querySelector('#switch-dark').addEventListener('click', function() {
	      document.querySelector('body').className = 'index black';
	    });

	    document.querySelector('#switch-light').addEventListener('click', function() {
	      document.querySelector('body').className = 'index white';
	    });
	  });
	})();
    `

	fmt.Println(printstr)
	filename := "./html/abcd_files/main.js"
	
	fi, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	 
	if err != nil {

		 panic(err)
	}
	 
	defer fi.Close()
	fi.WriteString(printstr)
	
}

func getkey(v []string, len int) string {
	key := "."
	for i := 0; i < len; i++ {
		key += v[i] + "."
	}
	return key

}

func Filter(result []*Tracker, deepth int) []*Tracker {

	filter := ""
	result2 := make([]*Tracker, 0)
	for i := range result {
		if result[i] == nil {
			continue
		}
		deepstr := strings.Split(result[i].Callstack, ".")
		if len(deepstr) == deepth+1 {

			if filter == "" {
				filter = getkey(deepstr, deepth)
				result2 = append(result2, result[i])
				result[i] = nil
				continue
			} else if filter == getkey(deepstr, deepth) {

				result2 = append(result2, result[i])
				result[i] = nil
				continue

			}

		}

	}
	return result2
}
