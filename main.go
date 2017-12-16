package main

import (
	Apollo "apollo4go/apollo"
	"time"
	"os"
	"fmt"
	"apollo4go/demo"
	"reflect"
)

func main() {
	os.Setenv("CONSUL_HTTP_ADDR", "127.0.0.1:8500")
	abc:=demo.DemoService{}.Hello
	def:=reflect.ValueOf(abc)
	params := make([]reflect.Value,1)
	params[0] = reflect.ValueOf("12")
	def.Call(params)
	Apollo.Register(demo.DemoService{})
	Apollo.Router("MyApollo_helloDaiIput",demo.DemoService{}.Hello)
	//Apollo.Run()
	for {
		var result string
		error := Apollo.Call("MyApollo_helloDaiIput", &result, "golang")
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println(result)
		time.Sleep(time.Second * 5)
	}
}
