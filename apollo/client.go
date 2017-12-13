package Apollo


import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "127.0.0.1:5197"
)

func Call(name string)  string{
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewApolloServiceClient(conn)

	r, err := c.Call(context.Background(), &Request{ServiceName: name})
	if err != nil {
		log.Fatal("could not call: %v", err)
	}
	log.Printf("Output is : %s %s", r.Message,r.Data)
	return r.Data
}