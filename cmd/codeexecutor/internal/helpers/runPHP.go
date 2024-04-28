package helpers

import (
	"context"
	// "fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunPHP(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/php-env:latest",
		Cmd:   []string{"/app/index.php"},
	}, &container.HostConfig{
		Binds: []string{tmpfilePath + ":/app/index.php"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	return resp, err
}
