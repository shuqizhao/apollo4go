package main

import (
	Apollo "apollo4go/apollo"
	"time"
	"os"
	"fmt"
)

func main(){
	os.Setenv("CONSUL_HTTP_ADDR", "127.0.0.1:8500")
	//Apollo.Run()
	for {
		var result string
		Apollo.Call("MyApollo_helloDaiIput",&result,"舒启钊")
		fmt.Println(result)
		time.Sleep(time.Second*5)
	}
}
