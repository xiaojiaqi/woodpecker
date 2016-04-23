package main

import (
	"fmt"
	"net"
	"time"
	"encoding/binary"
	"runtime"
)

func checkError(err error) (res bool) {

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func oneclient(service string, times int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

	if err != nil {
		return
	}

	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return
	}

	defer fmt.Println("conn close")
	defer conn.Close()

	checkError(err)

	msg := &Message{
		filetype:    111,
		serviceid:   222,
		timestamp:   uint32(time.Now().UnixNano()),
		messagetype: 333,
		bodylen:     10,
	} //msg init
	msg.body = make([]byte, msg.bodylen)
	var i uint32
	for i = 0; i < msg.bodylen; i++ {

		msg.body[i] = byte(i + 65)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(msg.timestamp))
	Send(conn, header, 4)
	binary.BigEndian.PutUint32(header, uint32(msg.messagetype))
	Send(conn, header, 4)

	binary.BigEndian.PutUint32(header, uint32(msg.bodylen))
	Send(conn, header, 4)

	Send(conn, msg.body, int(msg.bodylen))

}

func Send(conn *net.TCPConn, p []byte, buflen int) {
	for i := 0; i < buflen; i++ {

		_, err := conn.Write(p[i : i+1])
		fmt.Println("send ", p[i:i+1])
		time.Sleep(3 * time.Second)
		if err != nil {
			checkError(err)
			return
		}
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	service := "127.0.0.1:9955"
	oneclient(service, 1)

}

type Message struct {
	filetype    uint32 // 0 写meta
	serviceid   uint32 // id
	timestamp   uint32 //
	messagetype uint32
	bodylen     uint32
	body        []byte
}
