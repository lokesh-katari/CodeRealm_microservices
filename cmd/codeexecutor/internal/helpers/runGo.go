package helpers

import (
	"context"
	// "fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunGo(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/golang-env:latest",
		Cmd:   []string{"/app/main.go"},
	}, &container.HostConfig{
		Binds: []string{tmpfilePath + ":/app/main.go"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	return resp, err
}
