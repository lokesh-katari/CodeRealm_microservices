package codeservice

import (
	// "context"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lokesh-katari/code-realm/cmd/codeexecutor/internal/helpers"
	"os"

	// "github.com/docker/docker/api/types"
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

	// Create a map to hold the JSON response
	jsonResponse := map[string]string{
		"output": output,
		"error":  "",
		"lang":   code.Language,
	}
	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonResponseBytes))
	return string(jsonResponseBytes), nil
}
