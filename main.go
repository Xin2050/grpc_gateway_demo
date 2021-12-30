package main

import (
	"context"
	hello_world_pb_v1 "github.com/Xin2050/grpc_gateway_demo/proto/gen/hello_world/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
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

	// Create a gRPC server object gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	hello_world_pb_v1.RegisterGreeterServer(s, &server{})
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()
	client := hello_world_pb_v1.NewGreeterClient(conn)
	result, err := client.SayHello(context.Background(), &hello_world_pb_v1.HelloRequest{
		Name: "Leon",
	})
	if err != nil {
		log.Fatalf("Fail to create blog %v", err)
		return
	}
	log.Println(result.Message)

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = hello_world_pb_v1.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
