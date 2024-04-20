package codeservice

import (
	// "context"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"lokesh-katari/code-realm/cmd/codeexecutor/internal/helpers"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Code struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func CodeSubmission(resCode string, lang string) (string, error) {

	var tmpfilePath string
	var resp container.CreateResponse
	var code Code
	var tmpfile *os.File

	code.Language = lang
	code.Code = resCode

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	switch code.Language {
	case "javascript":
		tmpfile, err := ioutil.TempFile("", "example*.js")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunJs(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "python":
		tmpfile, err := ioutil.TempFile("", "example*.py")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunPy(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "java":
		fmt.Println("java code")
		tmpfile, err := ioutil.TempFile("", "example*.java")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunJava(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	}

	// Display the name of the temporary file.
	// fmt.Println(tmpfilePath, "this is  filepath")

	// Close the file.
	tmpfile.Close()
	// fmt.Println(resp.ID, "this is response id")
	defer os.Remove(tmpfilePath)

	// if err != nil {
	// 	panic(err)
	// }

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("container started successfully")

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if _, err := stdcopy.StdCopy(&buf, &buf, out); err != nil {
		panic(err)
	}
	cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	// Convert the buffer content to a string
	output := buf.String()

	// // Extract the runtime and status from the output
	// runtimeRegex := regexp.MustCompile(`Runtime: (\d+\.\d+) seconds`)
	// statusRegex := regexp.MustCompile(`All Test Cases Passed: (True|False|true|false)`)

	// var runtime float64
	// var status bool
	// var passedTestCases []int

	// if match := runtimeRegex.FindStringSubmatch(output); len(match) > 1 {
	// 	runtime, _ = strconv.ParseFloat(match[1], 64)
	// } else {
	// 	runtime = 0.0
	// }

	// if match := statusRegex.FindStringSubmatch(output); len(match) > 1 {
	// 	status = (match[1] == "True" || match[1] == "true")
	// }

	// // Extract the test cases that passed
	// testCaseRegex := regexp.MustCompile(`Test Case (\d+): .*Status: Accepted`)
	// matches := testCaseRegex.FindAllStringSubmatch(output, -1)
	// for _, match := range matches {
	// 	testCaseNumber, _ := strconv.Atoi(match[1])
	// 	passedTestCases = append(passedTestCases, testCaseNumber)
	// }

	// // Create a map to hold the JSON response
	// jsonResponse := map[string]interface{}{
	// 	"output":          output,
	// 	"runtime":         runtime,
	// 	"lang":            code.Language,
	// 	"status":          status,
	// 	"passedTestCases": passedTestCases,
	// }

	// jsonResponseBytes, err := json.Marshal(jsonResponse)
	// if err != nil {
	// 	panic(err)
	// }
	// return string(jsonResponseBytes), nil
	return output, nil
}
