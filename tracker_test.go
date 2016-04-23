package wood

import (
	"fmt"
	"testing"
	//"runtime"
	    . "github.com/woodpecker/tracker"
)

func Test_wood_Ack(te *testing.T) {
	Testcase ()
}

func Testcase () {
    L := &LocalSender{}

	Gsender = L

	t := &Tracker{}
	t.Init(111)
	defer t.Close()
	
	//fmt.Println(t.sid)     //1
	t.Start() //1
	//fmt.Println(t)

	
A(t.Clone()) //1.
A(t.Clone()) //1.
A(t.Clone()) //1.
	t.End()      //1
	//fmt.Println(t)

	t.Start() //2
	//fmt.Println(t)

	t.End() //2
	//fmt.Println(t)

	t.Start() //3
	//fmt.Println(t)

	t.End() //3
	//fmt.Println(t)
	 t.Metric("main func ", 1111123)

	fmt.Println("===============")
	for _, i := range L.Vector {
		fmt.Println(i)
	}
}

func A(t *Tracker) {

	t.NewSpace() //1.
    defer t.Close()
	t.Start() //1.1

	//fmt.Println(t)
	t.Event("A func ")
	

	t.End()      // 1.1
	B(t.Clone()) // 1.2

}

func B(t *Tracker) {

	t.NewSpace()
	defer t.Close()

	t.Start() //1.2

	//fmt.Println(t)

    t.Metric("A func ", 333)
	
	t.End()
	C(t.Clone())

}

func C(t *Tracker) {

	t.NewSpace()
defer t.Close()
	t.Start()

	//fmt.Println(t)

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
