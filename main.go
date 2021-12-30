package main

import (
	"context"
	hello_world_pb_v1 "github.com/Xin2050/grpc_gateway_demo/proto/gen/hello_world/v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	//hello_world_pb_v1.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}
func (s *server) SayHello(ctx context.Context, in *hello_world_pb_v1.HelloRequest) (*hello_world_pb_v1.HelloReply, error) {
	return &hello_world_pb_v1.HelloReply{Message: in.Name + " world"}, nil
}
func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	hello_world_pb_v1.RegisterGreeterServer(s, &server{})
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}
