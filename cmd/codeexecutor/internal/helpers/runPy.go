package helpers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunPy(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	fmt.Println("attempting to run cotainer")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "lokeshkatari/python-env:latest",
		Cmd:   []string{"/app/code.py"},
	}, &container.HostConfig{
		// AutoRemove: true,
		Binds: []string{tmpfilePath + ":/app/code.py"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("container created successfully", resp, "this is resp frm runpy")

	return resp, err
}
