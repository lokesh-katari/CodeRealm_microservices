package main

import (
	"context"
	"fmt"
	"log"
	"lokesh-katari/code-realm/cmd/client/codeExecutionpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":50052"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := codeExecutionpb.NewCodeExecutionServiceClient(conn)
	res, err := client.ExecuteCode(context.Background(), &codeExecutionpb.ExecuteCodeRequest{
		Language:  "python",
		Code:      "print('hello world this is lokesh')",
		InputData: []string{"hello world"},
	})
	fmt.Println(err, "thos is user login")
	fmt.Println(res, "thos is user login")

}
