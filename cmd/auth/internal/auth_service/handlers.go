package authservice

import (
	"context"
	"fmt"
	pb "lokesh-katari/code-realm/cmd/auth/internal/proto/auth"

	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	// pb.UnimplementedAuthServiceServer

	pb.AuthServiceServer
	authservice    AuthService
	UserRepository PostgresUserRepository
}

func NewServer(authservice AuthService) *Server {
	return &Server{
		authservice: authservice,
	}
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user, err := s.authservice.RegisterUser(ctx, req.Email, req.Password)
	fmt.Println("user", user, req, "this is request")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	return &pb.RegisterUserResponse{
		Token: user.Email,
		User: &pb.User{
			Email: user.Email,
			Name:  user.Username,
		},
	}, nil

}
func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := s.authservice.LoginUser(ctx, req.Email, req.Password)

	if err != nil {
		log.Println("error in the login user", err)
	}
	return &pb.LoginUserResponse{
		Token: user,
	}, nil
}

func (s *Server) LogoutUser(ctx context.Context, req *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	if req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Token is required")
	}

	if err := s.authservice.LogoutUser(ctx, req.Token); err != nil {
		return nil, status.Errorf(codes.Internal, "Error logging out user: %v", err)
	}
	return &pb.LogoutUserResponse{
		Success: true,
	}, nil

}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.authservice.GetUser(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting user: %v", err)
	}
	return &pb.GetUserResponse{
		User: &pb.User{
			Email: user.Email,
			Name:  user.Username,
		},
	}, nil
}
