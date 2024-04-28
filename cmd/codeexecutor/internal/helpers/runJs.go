package helpers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunJs(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/js-env:latest",
		Cmd:   []string{"/app/code.js"},
	}, &container.HostConfig{

		Binds: []string{tmpfilePath + ":/app/code.js"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	return resp, err
}
