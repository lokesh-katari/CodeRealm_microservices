package helpers

import (
	"context"
	// "fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunJava(ctx context.Context, cli *client.Client, tmpfilePath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "java-env",
		Cmd:   []string{"/app/Solution.java"},
	}, &container.HostConfig{
		AutoRemove: true,
		Binds:      []string{tmpfilePath + ":/app/Solution.java"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	// fmt.Println("Java container created", resp.ID)
	return resp, err
}
