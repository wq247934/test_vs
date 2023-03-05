package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"syscall"
)

func main() {
	fmt.Println(syscall.SIGHUP.String())
	return
	err := reloadNginxByContainer("")
	if err != nil {
		panic(err)
		return
	}
}

func reloadNginxByContainer(name string) error {
	name = "/" + name
	DockerCli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	DockerCli.NegotiateAPIVersion(context.Background())
	containers, err := DockerCli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	for _, container := range containers {
		for _, n := range container.Names {
			if name == n {
				return DockerCli.ContainerKill(context.Background(), container.ID, syscall.SIGHUP.String())
			}
		}
	}
	return err
}
