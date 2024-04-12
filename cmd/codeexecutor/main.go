package main

import (
	"log"
	"net"

	codeservice "lokesh-katari/code-realm/cmd/codeexecutor/internal/code_service"
	pb "lokesh-katari/code-realm/cmd/codeexecutor/internal/proto/codeExecutionpb"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}
	s := grpc.NewServer()
	server := codeservice.NewServer()
	pb.RegisterCodeExecutionServiceServer(s, server)

	reflection.Register(s)
	log.Println("gRPC server started at port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
