package tracker

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	/*
		"fmt"
		"math/rand"
		"runtime"
		"strconv"
		"strings"
		"time"
	*/
	"net"
	"time"
	//"errors"
	//"io"
	"encoding/binary"
	//. "github.com/woodpecker/track"
	"strconv"
)

const (
	ChannelSize         = 10000
	MessageHeaderLength = 12
	// TrackStart = 0
	//TrackEnd   = 1
	//TrackEvent  = 2
	//TrackMetric = 3
)

var channel chan *Tracker

 
func checkError(err error) (res bool) {

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type LocalSender struct {
	Vector []*Tracker
}

type RemoteSender struct {
	host string
	port int
	conn *net.TCPConn
}

func (v *LocalSender) Send(t *Tracker) {

	if v.Vector == nil {
		v.Vector = make([]*Tracker, 0)
	}
	v.Vector = append(v.Vector, t)

}

func (v *LocalSender) SendToServer() {

}

func (v *RemoteSender) Send(t *Tracker) {

	_ = PushtoQueue(&channel, t)
}

func (v *RemoteSender) SendToServer() {

	bufflen := 65 * 1024 * 100
	buff := make([]byte, bufflen)
	index := 0
	var timestamp uint32
	var serviceid uint32
	var messagetype uint32
	var datalen uint32

	header := make([]byte, 4)
	messageindex := 1

	for {
	    index = 0
		Messages := GetSomeMessageFromQueue(&channel, 1000)
		for i := range Messages {
			timestamp = uint32(Messages[i].Timestamp/ 1e9)
			serviceid   = uint32(Messages[i].Serviceid)
			messagetype = uint32(Messages[i].Trackertype)
			messagetype = uint32(messageindex)
			
			messageindex += 1
            fmt.Println(messageindex, Messages[i])
			out, err := proto.Marshal(Messages[i])
			if err != nil {
			    continue
			}
			datalen = uint32(len(out))

			if index+MessageHeaderLength+int(datalen) < bufflen {

				binary.BigEndian.PutUint32(header, timestamp)
				CopyData(header, buff  , 4, 0, index)
				index += 4
				
				binary.BigEndian.PutUint32(header, serviceid)
				CopyData(header, buff  , 4, 0, index)
				index += 4
				
				binary.BigEndian.PutUint32(header, messagetype)
				CopyData(header, buff  , 4, 0, index)
				index += 4
				binary.BigEndian.PutUint32(header, datalen)
				CopyData(header, buff  , 4, 0, index)
				index += 4
              
				CopyData(out,buff,   int(datalen), 0, index)
				index += int(datalen)
				
				
				

			} else {
				break
			}

		}
		Send(v.conn, buff, index)
	}
}

func CopyData(src []byte, dest []byte, l int, aindex int, bindex int) {
	for i := 0; i < l; i++ {
		dest[bindex+i] = src[aindex+i]

	}

}


func Send(conn *net.TCPConn, p []byte, buflen int) {
	for i := 0; i < buflen; i++ {
        
		_, err := conn.Write(p[i : i+1])
		//fmt.Println("send ", p[i : i + 1 ])
		//time.Sleep(3 * time.Second)
		if err != nil {
			checkError(err)
			return
		}

	}
	//fmt.Println("=====================")

}

func InitTrackerClient(host string, port int) {
	channel = make(chan *Tracker, ChannelSize)

	L := &RemoteSender{}
    L.host = host
    L.port = port
    Gsender = L
    
	go SendTrack(L)
	//bug fix me
	time.Sleep(1 * time.Second)
}

func SendTrack(L *RemoteSender) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", L.host+":"+strconv.Itoa(int(L.port)))

	if err != nil {
	    fmt.Print(err)
		return
	}

	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return
	}
    L.conn = conn
	

	L.SendToServer()
}

func InitTrackerServer(host string, port int) {
	channel = make(chan *Tracker, ChannelSize*10)

}

func RecvfromQueue(queue *chan *Tracker) (*Tracker, error) {

	select {
	case v := <-*queue:
		return v, nil
	default:
		//fmt.Println("it is empty")
		return nil, errors.New("it is full")
	}
}

func RecvfromQueueBlock(queue *chan *Tracker) (*Tracker, error) {

	v := <-*queue
	return v, nil
}

func GetSomeMessageFromQueue(queue *chan *Tracker, least int) []*Tracker {

	postList := make([]*Tracker, 0)

	v, _ := RecvfromQueueBlock(queue)
	postList = append(postList, v)
	for i := 1; i < least; i++ {

		v, err := RecvfromQueue(queue)
		if err != nil {
			return postList
		}
		postList = append(postList, v)
	}
	return postList
}

func PushtoQueue(queue *chan *Tracker, req *Tracker) error {

	select {
	case *queue <- req:

		return nil
	default:

		return errors.New("it is full")
	}
}

func GetIntValue(key string) int {

	result, err := strconv.Atoi(key)
	if err != nil {
		result = 0
	}
	return result
}
