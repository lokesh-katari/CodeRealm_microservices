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
	user, err := s.authservice.RegisterUser(ctx, req.Name, req.Email, req.Password)
	fmt.Println("user", user, req, "this is request")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	token, err := user.GenerateJWT()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.RegisterUserResponse{
		Token: token,
		User: &pb.User{
			Email: user.Email,
			Name:  user.Name,
		},
	}, status.Error(codes.OK, "User registered successfully")

}
func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	token, err := s.authservice.LoginUser(ctx, req.Email, req.Password)

	if err != nil {
		log.Println("error in the login user", err)
	}
	return &pb.LoginUserResponse{
		Token: token,
	}, status.Error(codes.OK, "User logged in successfully")
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
	}, status.Error(codes.OK, "User logged out successfully")

}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Println("this is get user request", req.Token)
	user, err := s.authservice.GetUser(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting user: %v", err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Email:               user.Email,
			Name:                user.Name,
			Easy_ProblemCount:   int32(user.Easy_Problem_count),
			Medium_ProblemCount: int32(user.Medium_Problem_count),
			Hard_ProblemCount:   int32(user.Hard_Problem_count),
			SolvedProblems:      user.Submission,
		},
	}, status.Error(codes.OK, "User fetched successfully")
}

func (s *Server) UpdateUserSubmissions(ctx context.Context, req *pb.UpdateUserSubmissionsRequest) (*pb.UpdateUserSubmissionsResponse, error) {
	fmt.Println("this is get user request from usersubmissions", req.Token)
	user, err := s.authservice.UpdateUserSubmissions(ctx, req.Token, req.Queid, req.Difficulty)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting user: %v", err)
	}

	return &pb.UpdateUserSubmissionsResponse{
		User: &pb.User{
			Email:               user.Email,
			Name:                user.Name,
			Easy_ProblemCount:   int32(user.Easy_Problem_count),
			Medium_ProblemCount: int32(user.Medium_Problem_count),
			Hard_ProblemCount:   int32(user.Hard_Problem_count),
			SolvedProblems:      user.Submission,
		},
	}, status.Errorf(codes.OK, "User updated successfully")
}

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {

	err := s.authservice.ChangePassword(ctx, req.Token, req.OldPassword, req.NewPassword)
	if err != nil {
		return &pb.ChangePasswordResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "Error getting user: %v", err)
	}

	return &pb.ChangePasswordResponse{
		Success: true,
	}, status.Errorf(codes.OK, "Password changed successfully")
}
