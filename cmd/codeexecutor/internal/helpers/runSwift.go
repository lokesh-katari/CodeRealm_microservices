package helpers

import (
	"context"
	// "fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunSwift(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/swift-env:latest",
		Cmd:   []string{"/app/Solution.swift"},
	}, &container.HostConfig{
		Binds: []string{tmpfilePath + ":/app/Solution.swift"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	return resp, err
}
