package main

import (
	"log"
	auth "lokesh-katari/code-realm/cmd/auth/internal/auth_service"
	pb "lokesh-katari/code-realm/cmd/auth/internal/proto/auth"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}
	log.Println("Server started at port 50051")
	userRepo, err := auth.NewPostgresUserRepository()
	jwtManager := auth.NewJWTManager("secret")

	authService := auth.NewAuthServiceImpl(userRepo, jwtManager)
	interceptor := auth.NewAuthInterceptor(jwtManager)
	server := auth.NewServer(authService)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
	)
	pb.RegisterAuthServiceServer(s, server)
	reflection.Register(s)
	log.Println("gRPC server started at port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
