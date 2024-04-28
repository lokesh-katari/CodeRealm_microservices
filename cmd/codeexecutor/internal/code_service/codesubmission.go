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
		fmt.Println("inside the js code")

		tmpfile, err := ioutil.TempFile("", "example*.js")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("tmpfile created successfully", tmpfile.Name(), "this is tmpfile name", code.Code, "this is code ")
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
	case "csharp":
		tmpfile, err := ioutil.TempFile("", "example*.cs")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunCsharp(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "cpp":
		tmpfile, err := ioutil.TempFile("", "example*.cpp")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunCpp(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "golang":
		tmpfile, err := ioutil.TempFile("", "example*.go")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunGo(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "ruby":
		tmpfile, err := ioutil.TempFile("", "example*.rb")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunRuby(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "swift":
		tmpfile, err := ioutil.TempFile("", "example*.swift")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunSwift(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "c":
		tmpfile, err := ioutil.TempFile("", "example*.c")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunC(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "rust":
		tmpfile, err := ioutil.TempFile("", "example*.rs")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunRust(ctx, cli, tmpfilePath)
		if err != nil {
			panic(err)
		}
	case "php":
		tmpfile, err := ioutil.TempFile("", "example*.php")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tmpfile, code.Code)
		tmpfilePath = tmpfile.Name()
		resp, err = helpers.RunPHP(ctx, cli, tmpfilePath)
		if err != nil {

			panic(err)
		}

	}

	fmt.Println(resp, "this is response from the container")
	// Close the file.
	tmpfile.Close()
	fmt.Println(tmpfilePath, "this is tmpfile path")
	// fmt.Println(resp.ID, "this is response id")
	// defer os.Remove(tmpfilePath)

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
	fmt.Println("container logs obtained successfully", out)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if _, err := stdcopy.StdCopy(&buf, &buf, out); err != nil {
		panic(err)
	}
	// defer cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	// Convert the buffer content to a string
	output := buf.String()

	return output, nil
}
