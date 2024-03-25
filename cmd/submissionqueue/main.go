package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type CodeExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "code-exec-requests",
		GroupID: "my-group",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		var codeExecutionRequest CodeExecutionRequest
		json.Unmarshal(msg.Value, &codeExecutionRequest)
		fmt.Println(codeExecutionRequest.Code, codeExecutionRequest.Language)

		fmt.Printf("Received message: %s\n", string(msg.Value))
		fmt.Println(string(msg.Key))
	}
}
