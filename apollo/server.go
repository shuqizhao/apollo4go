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

type Server struct {
	Port int
	IP string
}

func (s *Server) Call(ctx context.Context, in *Request) (*Response, error) {
	return &Response{Message: "Hello " + in.ServiceName}, nil
}

func Run()  {
	go func(){
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		RegisterApolloServiceServer(s, &Server{})
		s.Serve(lis)
	}()

}