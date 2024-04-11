package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"lokesh-katari/code-realm/cmd/client/codeExecutionpb"
	"lokesh-katari/code-realm/cmd/client/db"
	"lokesh-katari/code-realm/cmd/client/models"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":50052"
)

type CodeExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	PID      string `json:"pid"`
	ReqType  string `json:"reqType"`
	QueID    string `json:"queId"`
	Email    string `json:"email"`
}

type CodeExecutionResponse struct {
	Output string `json:"output"`
}

var REDIS_URI = "redis://default:vjIGMyBfPrVKyR1l7F12Gf0SxvHofMmq@redis-10614.c13.us-east-1-3.ec2.cloud.redislabs.com:10614"

func main() {
	conn, err := grpc.Dial(os.Getenv("GRPC_SERVER")+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer db.Client.Disconnect(context.TODO())
	opt, err := redis.ParseURL(REDIS_URI)
	if err != nil {
		fmt.Println(err)
	}

	rclient := redis.NewClient(opt)
	pong, err := rclient.Ping(context.Background()).Result()
	fmt.Println(pong, err)
	submissionreader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BOOTSTRAP_SERVERS")},
		Topic:   "code-submission-request",
		GroupID: "submission-group",
	})

	runreader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BOOTSTRAP_SERVERS")},
		Topic:   "code-run-request",
		GroupID: "run-group",
	})
	defer submissionreader.Close()
	defer runreader.Close()

	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	codeExecutionChannel := make(chan CodeExecutionRequest)

	go processMessages(submissionreader, codeExecutionChannel)
	go processMessages(runreader, codeExecutionChannel)

	for req := range codeExecutionChannel {
		fmt.Println("Received code execution request", req)
		go executeAndStore(rclient, conn, req)
	}
	defer rclient.Close()
}
func processMessages(reader *kafka.Reader, ch chan<- CodeExecutionRequest) {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var codeExecutionRequest CodeExecutionRequest
		err = json.Unmarshal(msg.Value, &codeExecutionRequest)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		ch <- codeExecutionRequest
	}
}

func executeAndStore(rclient *redis.Client, conn *grpc.ClientConn, req CodeExecutionRequest) {
	client := codeExecutionpb.NewCodeExecutionServiceClient(conn)

	var problem models.Problem
	// Declare the variable "problem"
	if req.ReqType == "submit" {

		opts := options.FindOne().SetProjection(bson.M{"templates." + req.Language: 1})
		err := db.CodeQueCollection.FindOne(context.TODO(), bson.M{"problemId": req.QueID}, opts).Decode(&problem)
		if err != nil {
			log.Printf("Error finding problem in MongoDB: %v", err)
			return
		}
		req.Code, err = GenerateCode(req.Language, req.Code, problem)
	}
	res, err := client.ExecuteCode(context.Background(), &codeExecutionpb.ExecuteCodeRequest{
		Language:  req.Language,
		Code:      req.Code,
		InputData: []string{"hello world"},
	})
	if err != nil {
		log.Printf("Error when calling ExecuteCode: %v", err)
		return
	}
	if req.ReqType == "submit" {
		fmt.Println("Storing submission in MongoDB", res, "this is res")
		_, err = db.SubmissionCollection.InsertOne(context.TODO(), models.CodeSubmission{
			PID:         req.PID,
			QueID:       req.QueID,
			Email:       req.Email,
			Language:    req.Language,
			Code:        req.Code,
			Output:      res.Output,
			SubmittedAT: time.Now(),
			Runtime:     "0",
			Memory:      "0",
		})
		if err != nil {
			log.Printf("Error storing submission in MongoDB: %v", err)
		}
		return // Do not store output in Redis
	}
	err = rclient.Set(context.Background(), req.PID, res.Output, 3*time.Minute).Err()
	fmt.Println("Stored output in Redis")

	if err != nil {
		log.Printf("Error storing output in Redis: %v", err)
	}
}
