package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	container_types "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/container"
	"os"
)

var (
	ctx = context.Background()
	tendermintImage = "tendermint/tendermint:latest"

	seedIP string
)

func main() {
	// get docker client
	cli, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println("couldnt set up client: ", err)
		os.Exit(1)
	}
	fmt.Println("made docker client", cli.ClientVersion())
	// spin up seed node
	_, err = spinUpSeedNode(cli)
	if err != nil {
		fmt.Println("couldn't spin up seed node: ", err)
		os.Exit(1)
	}
	fmt.Println("spun up seed node?")
	// get peerID@IP:PORT of seed node
	// set it as ENV VAR
	// feed it to other clients
}

func spinUpSeedNode(cli *client.Client) (*container.Container, error) {
	_, err := cli.ImageBuild(ctx, nil, types.ImageBuildOptions{
		Tags: []string{"seed"},
		Dockerfile: tendermintImage,

	})
	if err != nil {
		return nil, err
	}

	_, err = cli.ContainerCreate(ctx, &container_types.Config{
		Image: "seed",
	}, nil, nil, nil, "seed")
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		fmt.Println(container.Names)
	}
	return nil, nil
}


