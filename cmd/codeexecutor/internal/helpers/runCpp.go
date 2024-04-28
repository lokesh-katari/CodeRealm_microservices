package helpers

import (
	"context"
	// "fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunCpp(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/gcc-env:latest",
		Cmd:   []string{"/app/Solution.cpp"},
	}, &container.HostConfig{
		Binds: []string{tmpfilePath + ":/app/Solution.cpp"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	return resp, err
}
