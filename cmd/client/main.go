package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"lokesh-katari/code-realm/cmd/client/codeExecutionpb"
	"lokesh-katari/code-realm/cmd/client/db"
	"lokesh-katari/code-realm/cmd/client/models"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CodeExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	PID      string `json:"pid"`
	ReqType  string `json:"reqType"`
	QueID    string `json:"queId"`
	Email    string `json:"email"`
	Title    string `json:"title"`
}

type CodeExecutionResponse struct {
	Output string `json:"output"`
}

var REDIS_URI = os.Getenv("REDIS_URI")

func main() {
	err := godotenv.Load()
	if err != nil {
		// Handle error loading .env file
		log.Println("Error loading .env file")
	}

	conn, err := grpc.Dial(os.Getenv("GRPC_URI_CODE_CLIENT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	fmt.Println("Connected to gRPC", conn, "this is connection", os.Getenv("GRPC_URI_CODE_CLIENT"))
	defer db.Client.Disconnect(context.TODO())
	opt, err := redis.ParseURL(os.Getenv("REDIS_URI"))
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
	fmt.Println("Connected to Kafka")

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
	fmt.Println("Processing messages")
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

func executeAndStore(rclient *redis.Client, conn *grpc.ClientConn, req CodeExecutionRequest) error {
	client := codeExecutionpb.NewCodeExecutionServiceClient(conn)
	fmt.Println("Executing code", req, "this is from the execute and store", client, conn)

	var Template models.Templates
	if req.ReqType == "submit" || req.ReqType == "run" {
		fmt.Println("Request type", req.ReqType, req.QueID, "this is queId")

		// templateId ,_ := primitive.ObjectIDFromHex(CodeQue.TemplateID)
		// Declare the variable "problem"
		problemId, _ := primitive.ObjectIDFromHex(req.QueID)

		var CodeQue models.CodeQue

		err := db.CodeQueCollection.FindOne(context.TODO(), bson.M{"_id": problemId}).Decode(&CodeQue)
		if err != nil {
			log.Printf("Error finding problem in MongoDB: %v", err)
			return err
		}
		err = db.TemplateCollection.FindOne(context.TODO(), bson.M{"_id": CodeQue.TemplateID}).Decode(&Template)
		if err != nil {
			log.Printf("Error finding template in MongoDB: %v", err)
			return err
		}
		// fmt.Println("Template", Template)
		req.Code, err = GenerateCode(req.Language, req.Code, Template, req.ReqType)
	}
	fmt.Println("Generated code", req.Code, req.Language)
	res, err := client.ExecuteCode(context.Background(), &codeExecutionpb.ExecuteCodeRequest{
		Language: req.Language,
		Code:     req.Code,
	})
	fmt.Println("Executed code", res)
	if err != nil {
		log.Printf("Error when calling ExecuteCode: %v", err)
		return err
	}
	if req.ReqType == "submit" || req.ReqType == "run" {
		// Extract the runtime and status from the output
		fmt.Println("Extracting output", "this is from the submit code client")

		output, err := Seperateoutput(res.Output, req.Language)
		if err != nil {
			log.Printf("Error separating output: %v", err)
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			// Handle error
			log.Printf("Error unmarshaling JSON: %v", err)
		}
		fmt.Println("Data", data, "this is from the submit code client")
		// passedCases, ok := data["passedTestCases"]
		passedTestCasesInterface := data["testcases"].([]interface{})

		// Convert each element to int
		var passedTestCases []int
		for _, v := range passedTestCasesInterface {
			passedTestCases = append(passedTestCases, int(v.(float64)))
		}

		if req.ReqType == "submit" {
			fmt.Println("inside the sunit client", passedTestCases)
			err = InsertSubmissionsAndUpdateCodeQue(req.QueID, res.Output, req, true, data["runtime"].(string), passedTestCases)
			if err != nil {
				log.Printf("Error inserting submission: %v", err)
			}
			return nil

		}

		err = rclient.Set(context.Background(), req.PID, output, 3*time.Minute).Err()
		if err != nil {
			log.Printf("Error storing output in Redis: %v", err)
		}
		return nil

	}
	err = rclient.Set(context.Background(), req.PID, res.Output, 3*time.Minute).Err()
	fmt.Println("Stored output in Redis")

	if err != nil {
		log.Printf("Error storing output in Redis: %v", err)
	}
	return nil
}

func InsertSubmissionsAndUpdateCodeQue(queId string, output string, req CodeExecutionRequest, status bool, runtime string, testcases []int) error {
	fmt.Println("Inserting submission", queId, output, req, status, runtime, testcases, "this is from the insert submission")
	fmt.Printf("%T,%T,%T,%T,%T,%T", queId, output, req, status, runtime, testcases)
	queID, err := primitive.ObjectIDFromHex(queId)
	if err != nil {
		return err
	}
	_, err = db.SubmissionCollection.InsertOne(context.TODO(), models.CodeSubmission{
		Title:       req.Title,
		Accepted:    status,
		PID:         req.PID,
		QueID:       queID,
		Email:       req.Email,
		Language:    req.Language,
		Code:        req.Code,
		Output:      output,
		SubmittedAT: time.Now(),
		Runtime:     runtime,
		Testcases:   testcases,
	})
	if err != nil {
		return err
	}
	UpdateCodeQue(queID, status)
	return nil
}

func UpdateCodeQue(queID primitive.ObjectID, status bool) error {

	var updateField string
	if status {
		updateField = "submissions.correct"
	} else {
		updateField = "submissions.wrong"
	}

	update := bson.M{"$inc": bson.M{updateField: 1}}

	// UpdateOne updates a single document matching the filter.
	result, err := db.CodeQueCollection.UpdateByID(context.Background(), queID, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("no document updated")
	}
	return nil
}

func GenerateCode(language string, userCode string, template models.Templates, reqType string) (string, error) {
	var finalCode string
	fmt.Println("Generating code")
	// fmt.Println(template)

	if reqType == "submit" {
		switch language {
		case "python":
			finalCode = userCode + "\n" + template.Python.HiddenTestCode
		case "javascript":
			finalCode = userCode + "\n" + template.JavaScript.HiddenTestCode
		case "golang":
			finalCode = userCode + template.Golang.HiddenTestCode
		case "java":
			finalCode = userCode + template.Java.HiddenTestCode
		case "c":
			finalCode = userCode + template.C.HiddenTestCode
		case "cpp":
			finalCode = userCode + template.Cpp.HiddenTestCode
		default:
			return "", fmt.Errorf("Invalid language: %s", language)
		}
	} else {
		switch language {
		case "python":
			finalCode = userCode + "\n" + template.Python.RunTestCode
		case "javascript":
			finalCode = userCode + "\n" + template.JavaScript.RunTestCode
		case "golang":
			finalCode = userCode + template.Golang.RunTestCode
		case "java":
			finalCode = userCode + template.Java.RunTestCode
		case "c":
			finalCode = userCode + template.C.RunTestCode
		case "cpp":
			finalCode = userCode + template.Cpp.RunTestCode
		default:
			return "", fmt.Errorf("Invalid language: %s", language)
		}
	}
	// finalCode = userCode + template.Python.RunTestCode
	fmt.Println("Generated code adsf", finalCode, reqType, "thisis requestype")
	return finalCode, nil
}

func Seperateoutput(output string, language string) (string, error) {
	runtimeRegex := regexp.MustCompile(`Runtime: (\d+\.\d+) seconds`)
	statusRegex := regexp.MustCompile(`All Test Cases Passed: (True|False|true|false)`)

	var runtime float64
	var status bool
	var passedTestCasesInt []int

	if match := runtimeRegex.FindStringSubmatch(output); len(match) > 1 {
		runtime, _ = strconv.ParseFloat(match[1], 64)
	} else {
		runtime = 0.0
	}

	if match := statusRegex.FindStringSubmatch(output); len(match) > 1 {
		status = (match[1] == "True" || match[1] == "true")
	}

	pattern := `Passed Test Cases:\s*\[([^]]+)\]`

	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Find the matches in the output string
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		// Extract the matched string containing the passed test cases array
		passedTestCasesStr := matches[1]

		// Split the string into individual test case values
		passedTestCases := regexp.MustCompile(`\s*,\s*`).Split(passedTestCasesStr, -1)
		for _, testCase := range passedTestCases {
			trimmedValue := strings.TrimSpace(testCase)
			tc, err := strconv.Atoi(trimmedValue)
			if err != nil {
				fmt.Println("Error converting test case value to integer:", err)
				return "", err
			}
			passedTestCasesInt = append(passedTestCasesInt, tc)
		}

		fmt.Println("Passed Test Cases:", passedTestCasesInt)
	} else {
		fmt.Println("Passed Test Cases not found in the output string.")
	}
	stringRuntime := strconv.FormatFloat(runtime, 'f', -1, 64)

	// Create a map to hold the JSON response
	jsonResponse := map[string]interface{}{
		"output":    output,
		"runtime":   stringRuntime,
		"language":  language,
		"status":    status,
		"testcases": passedTestCasesInt,
	}

	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		panic(err)
	}
	return string(jsonResponseBytes), nil
}
