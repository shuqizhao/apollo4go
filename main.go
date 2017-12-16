package main

import (
	Apollo "apollo4go/apollo"
	"time"
	"os"
	"fmt"
	"apollo4go/demo"
)

func main(){
	os.Setenv("CONSUL_HTTP_ADDR", "127.0.0.1:8500")
	Apollo.Register(&demo.DemoService{})
	//Apollo.Run()
	for {
		var result string
		error:=Apollo.Call("MyApollo_Hello",&result,"golang")
		if error!=nil{
			fmt.Println(error)
		}
		fmt.Println(result)
		time.Sleep(time.Second*5)
	}
}
