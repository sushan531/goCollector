package handlers

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var LoadHandlers []string

func GetDockerClient() (*client.Client, context.Context) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	return cli, ctx
}

func createContainer(cli *client.Client, imageName string, ctx context.Context) string {
	labels := make(map[string]string)
	labels["handler"] = "handler"
	resp, _ := cli.ContainerCreate(
		ctx,
		&container.Config{Image: imageName, Cmd: []string{"./loadHandler", "."}, Labels: labels},
		&container.HostConfig{
			Binds: []string{
				"/Users/spt/GolandProjects/collector/manager/loadHandler/config.yml:/usr/src/app/config.yml",
			},
		},
		nil, nil,
		"",
	)
	return resp.ID
}

func startContainer(cli *client.Client, containerId string, ctx context.Context) {
	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func KillDockerContainer(cli *client.Client, containerId string, ctx context.Context) {
	if err := cli.ContainerKill(ctx, containerId, "SIGKILL"); err != nil {
		panic(err)
	}
}

func RemoveDockerContainers(cli *client.Client, containerIds []string, ctx context.Context) {
	for _, containerId := range containerIds {
		_ = cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	}
}

func SpawnLoadHandlers(minHandlers int, cli *client.Client, ctx context.Context) []string {
	imageName := "loadhandler-handler"
	for i := 0; i < minHandlers; i++ {
		containerId := createContainer(cli, imageName, ctx)
		startContainer(cli, containerId, ctx)
		LoadHandlers = append(LoadHandlers, containerId)
	}
	for _, value := range LoadHandlers {
		fmt.Println("Container Id : ", value)
	}
	return LoadHandlers
}
