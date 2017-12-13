package Apollo

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {}

func (s *server) Call(ctx context.Context, in *Request) (*Response, error) {
	return &Response{Message: "Hello " + in.ServiceName}, nil
}

func Run()  {
	go func(){
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterApolloServiceServer(s, &server{})
		s.Serve(lis)
	}()

}