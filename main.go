package main

import "fmt"
import (
	Apollo "apollo4go/apollo"
	"time"
)

func main(){
	fmt.Println("123")
	//Apollo.Run()
	for {
		Apollo.Call("Apollo.DemoInterface.IValues$Apollo.DemoInterface.IValues_Hello")
		time.Sleep(time.Second*5)
	}
}
