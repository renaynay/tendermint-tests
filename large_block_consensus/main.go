package main

import (
	"context"
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	ctx             = context.Background()
	tendermintImage = "tendermint/tendermint:latest"

	seedIP string
)

func main() {
	// get docker client
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		fmt.Println("couldnt set up client: ", err)
		os.Exit(1)
	}
	fmt.Println("made docker client: ", client.Endpoint())
	// spin up seed node
	seed, err := spinUpSeedNode(client)
	if err != nil {
		fmt.Println("couldn't spin up seed node: ", err)
		os.Exit(1)
	}
	fmt.Println("spun up seed node!!! Seed ID: ", seed.ID)
	// get peerID@IP:PORT of seed node
	// set it as ENV VAR
	// feed it to other clients
}

func spinUpSeedNode(client *docker.Client) (*docker.Container, error) {
	seedContainer, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "seed",
		Config: &docker.Config{
			Image: tendermintImage,
			User: "root",
		},
		HostConfig: &docker.HostConfig{
			Mounts: []docker.HostMount{
				{
					Source: "./tendermint-seed/init/docker-entrypoint.sh",
					Target: "/usr/local/bin/docker-entrypoint.sh",
				},
				{
					Source: "./tendermint-seed/init/node_key.json",
					Target: "~/config/node_key.json",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("created container!!!!!!!: ", seedContainer.Created)
	return seedContainer, nil
}
