package codeservice

import (
	"context"
	"fmt"

	// "fmt"
	pb "lokesh-katari/code-realm/cmd/codeexecutor/internal/proto/codeExecutionpb"
)

type Server struct {
	pb.CodeExecutionServiceServer
	// codeservice CodeService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ExecuteCode(ctx context.Context, req *pb.ExecuteCodeRequest) (*pb.ExecuteCodeResponse, error) {
	// return s.codeservice.ExecuteCode(ctx, req)
	fmt.Println("Code received to the handler", req.Code)
	op, err := CodeSubmission(req.Code, req.Language)
	if err != nil {
		fmt.Println("Error in code submission", err)
	}
	return &pb.ExecuteCodeResponse{
		Output: op,
	}, nil
}
