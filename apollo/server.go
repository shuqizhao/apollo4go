package Apollo

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
	"reflect"
	"errors"
	consulapi "github.com/hashicorp/consul/api"
	"os"
	"strconv"
	"time"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
	"io"
	"encoding/base64"
	crand "crypto/rand"
	"strings"
	"encoding/json"
)

func init()  {
}

var Services []interface{}

var MethodValues = make(map[string] reflect.Value)
var MethodTypes = make(map[string] reflect.Method)

func Register(service interface{}) error {
	Services=append(Services,service)
	//object:=reflect.ValueOf(service)

//	objRef:=object.Elem()
	objType:=reflect.TypeOf(service).Elem()
	objValue:=reflect.ValueOf(service).Elem()
	serviceName:=""
	port:=""
	fieldCount:=objType.NumField()
	for i:=0;i<fieldCount ;i++  {
		field:=objType.Field(i)
		if field.Name=="ServiceMeta"{
			serviceName=field.Tag.Get("ServiceName")
			port=field.Tag.Get("ServicePort")
			break
		}
	}

	if serviceName=="" {
		return errors.New("Can't found ServiceName")
	}

	methodCount:=objType.NumMethod()
	for i:=0;i<methodCount ;i++  {
		methodType:=objType.Method(i)
		methodValue:=objValue.Method(i)
		MethodTypes[serviceName+"_"+methodType.Name]=methodType
		MethodValues[serviceName+"_"+methodType.Name]=methodValue
		fmt.Println(methodType.Name,methodType.Type)
	}
	BuildService(serviceName,port)
	return nil
}

type Server struct {
	Port int
	IP string
}

func (s *Server) Call(ctx context.Context, in *Request) (*Response, error) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknow panic")
			}
			log.Println(err)
		}
	}()
	if _, ok := MethodValues[in.ServiceName]; ok {
		methodValue := MethodValues[in.ServiceName]
		values := strings.Split(in.Data, "兲")
		values = values[0:len(values)-1]
		params := make([]reflect.Value, len(values))
		for i := 0; i < len(values); i++ {
			params[i] = reflect.ValueOf(values[i])
		}
		result := methodValue.Call(params)
		p1 := result[0].String()
		data := ""
		if len(result) > 0 {
			bytes, _ := json.Marshal(p1)
			data = string(bytes)
		}
		return &Response{Data: data, Code: "200"}, nil
	} else {
		message := "Can't found method " + in.ServiceName
		return &Response{Message: message, Code: "500"}, errors.New(message)
	}
}

func BuildService(serviceName string,port string) error {
	if port==""{
		port=GetPort(serviceName)
	}
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	//创建一个新服务。
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = serviceName+GetGuid()
	registration.Name = serviceName
	a, _ := strconv.Atoi(port)
	registration.Port = a
	registration.Tags = []string{GetHostName(), "golang"}
	registration.Address = GetLocalHostIp()

	//增加check。
	check := new(consulapi.AgentServiceCheck)
	checkPort:=GetPort("")
	check.TCP = GetLocalHostIp()+":"+checkPort
	//设置超时 5s。
	check.Timeout = "5s"
	//设置间隔 5s。
	check.Interval = "5s"
	//注册check服务。
	registration.Check = check
	log.Println("get check.HTTP:", check)

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Fatal("register server error : ", err)
	}
	go func(port string){
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal("failed to listen: %v", err)
		}

		for {
			conn, err := lis.Accept()
			if err != nil {
				continue
			}
			daytime := time.Now().String()
			conn.Write([]byte(daytime)) // don't care about return value
			conn.Close()                // we're finished with this client
		}
	}(checkPort)
	go func(port string) {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterApolloServiceServer(s, &Server{})
		s.Serve(lis)
	}(port)

	return nil
}

func GetLocalHostIp() string{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
	}
	ip:=""
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func GetHostName() string  {
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("%s", err)
	}
	return host
}

func GetPort(name string)  string{
	if name==""{
		port:=RandInt(2000,4000)
		for  {
			address := fmt.Sprintf(":%d",port)
			lis, err := net.Listen("tcp", address)
			if err == nil {
				defer lis.Close()
				break
			}
		}
		return strconv.Itoa(port)
	}else{
		md5Str := GetMd5String(name)
		first:=md5Str[0]
		last:=md5Str[len(md5Str)-1]
		a:=int(first)
		b:=int(last)
		port:=fmt.Sprintf("%d%d",a,b)
		address := fmt.Sprintf(":%s",port)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println(err)
			return GetPort("")
		}else {
			defer lis.Close()
			return port
		}
	}
	return ""
}

func RandInt(min,max int) int{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	port:=r.Intn(max)
	if port<min{
		RandInt(min,max)
	}
	return port
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}