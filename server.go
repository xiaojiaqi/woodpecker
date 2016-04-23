package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	. "github.com/woodpecker/tracker"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	MessageChannelSize = 100000
	RecordBufLen       = 1024000
	LogLiveTime        = 60
)

func checkError(err error, info string) (res bool) {

	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

var Mchannel chan *Message
var Writechannel chan *RecordBuf

func ReadByte(conn *net.TCPConn, buff []byte, needread int) (r error, b bool) {

	readed := 0
	for {
		n, err := conn.Read(buff[readed:needread])
		if err != nil {
			fmt.Println("  error:", err)
			if err != io.EOF {
				fmt.Println("read error:", err)
				return err, false
			} else {
				if readed == needread {
					return nil, true
				} else {
					return err, false
				}
			}
		}
		readed += n
		if readed == needread {
			return nil, true
		}
	}
}

func Handler(conn *net.TCPConn) {
	defer conn.Close()
	conn.SetLinger(-1)
	 

	var timestamp []byte
	var messagetypes []byte
	var messagelen []byte
	var serviceid []byte

	timestamp = make([]byte, 4)
	messagetypes = make([]byte, 4)
	messagelen = make([]byte, 4)
	serviceid = make([]byte, 4)

	for {
		if !ReadFromClient(conn, timestamp, serviceid, messagetypes, messagelen) {
			break
		}
	}

	fmt.Println("out ot ReadFromClient")
}

func ReadFromClient(conn *net.TCPConn, timestamp []byte, serviceid []byte, messagetypes []byte, messagelen []byte) bool {

	r, b := ReadByte(conn, timestamp, 4)
	if b != true {
		fmt.Println(r)
		return false
	}

	r, b = ReadByte(conn, serviceid, 4)
	if b != true {
		fmt.Println(r)
		return false
	}
	r, b = ReadByte(conn, messagetypes, 4)
	if b != true {
		fmt.Println(r)
		return false
	}
	r, b = ReadByte(conn, messagelen, 4)
	if b != true {
		fmt.Println(r)
		return false
	}

	itimestamp := binary.BigEndian.Uint32(timestamp[0:4])
	iserviceid := binary.BigEndian.Uint32(serviceid[0:4])
	imessagetype := binary.BigEndian.Uint32(messagetypes[0:4])
	ilen := binary.BigEndian.Uint32(messagelen[0:4])
	tm := &Message{}
	tm.timestamp = itimestamp
	tm.serviceid = iserviceid
	tm.body = make([]byte, ilen+16)

	fmt.Println(" itimestamp ", itimestamp, "serviceid :", iserviceid, " imessagetype", imessagetype, "  messagelen", ilen)

	r, b = ReadByte(conn, tm.body[16:], int(ilen))
	if b != true {
		fmt.Println(r)
		return false
	} else {
		tr := &Tracker{}
		err := proto.Unmarshal(tm.body[16:], tr)
		if err != nil {
			fmt.Println(err)
			return false
		} else {

			CopyData(timestamp, tm.body, 4, 0, 0)
			CopyData(serviceid, tm.body, 4, 0, 4)
			CopyData(messagetypes, tm.body, 4, 0, 8)
			CopyData(messagelen, tm.body, 4, 0, 12)
			tm.bodylen = ilen + 16
			fmt.Println(tr)
			
			fmt.Println(" this length ", ilen + 16 )

			MPushtoQueue(&Mchannel, tm)
		}

	}
	return true

}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	service := ":9955"
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", service)
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("listen failed", err)
		return
	}
	Mchannel = make(chan *Message, MessageChannelSize)
	Writechannel = make(chan *RecordBuf, MessageChannelSize/1000)
	totalBuffMap.Maps = make ( map[uint32]*BuffMap)
	go WriteThread()
	go Core()

	for {
		conn, _ := l.AcceptTCP()
		fmt.Println("local  ", conn.RemoteAddr().Network(), " ", conn.RemoteAddr().String())
		go Handler(conn)
	}

}

type Message struct {
	filetype    uint32 // 0 写meta
	serviceid   uint32 // id
	timestamp   uint32 //
	messagetype uint32
	bodylen     uint32
	body        []byte
}

func WriteThread() {

	for {
		r := <- Writechannel

		filename := strconv.Itoa(int(r.serviceid)) + "_" + strconv.Itoa(int(r.timestamp)/LogLiveTime) + ".data"
        WriteFile(filename, r)
	}

}

func WriteFile(filename string, r *RecordBuf) {

	fi, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0666)
	 
	if err != nil {

		 panic(err)
	}
	 
	defer fi.Close()
	fi.Seek(0, 2)
	info, err := fi.Stat()
	info.Size()
	header := make([]byte, 4)
	
	binary.BigEndian.PutUint32(header, r.timestamp)  // 4 btye timestamp
	CopyData(header, r.buff  , 4, 0, 0)              // 
	
	binary.BigEndian.PutUint32(header, r.WriteIndex - 8)
	CopyData(header, r.buff  , 4, 0, 4)
	fmt.Println("write Index ", r.WriteIndex)
	fmt.Println("file leng", r.WriteIndex )

	fi.Write(r.buff[0: r.WriteIndex  ])
	
	r = nil
	

}

type RecordBuf struct {
	WriteIndex uint32
	BuffLen    uint32
	serviceid  uint32 // id
	timestamp  uint32 //
	buff       []byte
}

type BuffMap struct {
	Buffmap map[uint32]*RecordBuf
}

type TotalBuffMap struct {
	Maps map[uint32]*BuffMap
}

var totalBuffMap TotalBuffMap

func (T *TotalBuffMap) AddMessage(r *Message) {

	var record *RecordBuf
	bufmap, ok := T.Maps[r.serviceid]
	if !ok {
		// not such map

		bufmap = &BuffMap{}
        bufmap.Buffmap = make (map[uint32]*RecordBuf)
		T.Maps[r.serviceid] = bufmap
	}

	record, ok = bufmap.Buffmap[r.timestamp]

	if !ok {
		record = &RecordBuf{}
		record.BuffLen = RecordBufLen
		record.WriteIndex = 8
		record.serviceid = r.serviceid
		record.timestamp = r.timestamp
		record.buff = make([]byte, RecordBufLen)
		bufmap.Buffmap[r.timestamp] = record
	}
	if r.bodylen+record.WriteIndex > record.BuffLen {
		Writechannel <- record
		record = &RecordBuf{}
		record.BuffLen = RecordBufLen
		record.WriteIndex = 8
		record.serviceid = r.serviceid
		record.timestamp = r.timestamp

		record.buff = make([]byte, RecordBufLen)
		bufmap.Buffmap[r.timestamp] = record
	}

	CopyData(r.body, record.buff, int(r.bodylen), 0, int(record.WriteIndex))
	
	record.WriteIndex += uint32(r.bodylen)
fmt.Println("Copy ", r.bodylen, " WriteIndex" ,record.WriteIndex)
}

func (T *TotalBuffMap) AddMessageL(l []*Message) {
	for i := range l {
		T.AddMessage(l[i])
	}
}

func (T *TotalBuffMap) fresh() {
    fmt.Println("totalBuffMap.fresh()")
	tnow := uint32(time.Now().Unix())
	for _, value := range T.Maps {
		for timestamp, rc := range value.Buffmap {
            fmt.Println(timestamp, tnow)
			if timestamp+LogLiveTime < tnow {
                fmt.Println("Writechannel <- rc")
				Writechannel <- rc
				delete(value.Buffmap, timestamp)
			}
		}
	}

	for key, value := range T.Maps {
		if len(value.Buffmap) == 0 {
			delete(T.Maps, key)
		}

	}

}

func Core() {
	var r *Message

	for {
		select {
		case r = <-Mchannel:
			//
			totalBuffMap.AddMessage(r)
			for {
				if len(Mchannel) > 0 {
					l := MGetSomeMessageFromQueue(&Mchannel, 1000)
					totalBuffMap.AddMessageL(l)
				} else {
					break
				}
			}
			totalBuffMap.fresh()
		case <-time.After(5 * time.Second): //超时5s
			//
            
			totalBuffMap.fresh()
		}
	}

}

// copy & paste
func MRecvfromQueue(queue *chan *Message) (*Message, error) {
	select {
	case v := <-*queue:
		return v, nil
	default:
		//fmt.Println("it is empty")
		return nil, errors.New("it is full")
	}
}

func MRecvfromQueueBlock(queue *chan *Message) (*Message, error) {

	v := <-*queue
	return v, nil
}

func MGetSomeMessageFromQueue(queue *chan *Message, least int) []*Message {

	postList := make([]*Message, 0)

	v, _ := MRecvfromQueueBlock(queue)
	postList = append(postList, v)
	for i := 1; i < least; i++ {

		v, err := MRecvfromQueue(queue)
		if err != nil {
			return postList
		}
		postList = append(postList, v)
	}
	return postList
}

func MPushtoQueue(queue *chan *Message, req *Message) error {

	select {
	case *queue <- req:

		return nil
	default:

		return errors.New("it is full")
	}
}
