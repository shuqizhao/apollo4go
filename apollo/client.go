package Apollo


import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
	"strings"
	"encoding/json"
	"errors"
)

func Call(name string,result interface{},args ...interface{})  error {
	index := strings.LastIndex(name, "_")
	serviceName := name[0:index]
	service, err := GetServer(serviceName)
	if err != nil {
		return nil
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.IP, service.Port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := NewApolloServiceClient(conn)
	jsonStr := ""
	for _, v := range args {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		jsonStr += string(jsonBytes) + "兲"
	}

	r, err := c.Call(context.Background(), &Request{ServiceName: name, Data: jsonStr})
	if err != nil {
		return err
	}
	//log.Printf("Output is : %s %s", r.Message, r.Data)
	if r.Code == "200" {
		json.Unmarshal([]byte(r.Data), &result)
		return nil
	} else {
		return errors.New(r.Message)
	}
}