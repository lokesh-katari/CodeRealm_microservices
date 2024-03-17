package main

import (
	"log"
	"net"

	auth "lokesh-katari/code-realm/cmd/auth/internal/auth_service"
	pb "lokesh-katari/code-realm/cmd/auth/internal/proto/auth"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}

	userRepo, err := auth.NewPostgresUserRepository()

	// Create service instance with the repository
	authService := auth.NewAuthServiceImpl(userRepo)
	server := auth.NewServer(authService)
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, server)
	reflection.Register(s)
	log.Println("gRPC server started at port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
