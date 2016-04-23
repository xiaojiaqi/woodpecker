package main

import (
	"time"
	. "github.com/woodpecker/tracker"
	"os"
	"runtime"
	 
)

func oneclient(service string, port int, serviceid int) {

	for i:= 0 ; i < 1; i++{
		t := &Tracker{}
		t.Init(serviceid)
		 
		t.Start() //1

		A(t.Clone()) //1.
	 
		t.End()      //1
	  
	}

}

 

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	service := "127.0.0.1"
	port := 9955
	InitTrackerClient(service, port)

	for i := 0; i < 1; i++ {
		oneclient(service, port, GetIntValue(os.Args[1]))
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func A(t *Tracker) {

	t.NewSpace() //1.

	t.Start() //1.1

	t.Event("A func ")

	t.End()      // 1.1
	B(t.Clone()) // 1.2

}

func B(t *Tracker) {

	t.NewSpace()

	t.Start() //1.2

	t.Metric("A func ", 333)

	t.End()
	C(t.Clone())

}

func C(t *Tracker) {

	t.NewSpace()

	t.Start()

	t.End()
	t.Metric("C func ", 123)
}

func A1(*Tracker) {

}

func B1(*Tracker) {

}

func C1(*Tracker) {

}

func A2(*Tracker) {

}

func B2(*Tracker) {

}

func C2(*Tracker) {

}
