package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	fmt.Println("killing container....")
	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: seed.ID})
	if err != nil {
		fmt.Println("err removing container: ", err)
		os.Exit(1)
	}
	// get peerID@IP:PORT of seed node
	// set it as ENV VAR
	// feed it to other clients
}

func spinUpSeedNode(client *docker.Client) (*docker.Container, error) {
	entrypointPath, err := filepath.Abs("./tendermint-seed/init/docker-entrypoint.sh")
	if err != nil {
		return nil, err
	}
	nodekeyPath, err := filepath.Abs("./tendermint-seed/init/node_key.json")
	if err != nil {
		return nil, err
	}

	seedContainer, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "seed",
		Config: &docker.Config{
			Image: tendermintImage,
			User: "root",
			Entrypoint: []string{"chmod u+x /usr/local/bin/docker-entrypoint.sh"},
		},
		HostConfig: &docker.HostConfig{
			Mounts: []docker.HostMount{
				{
					Source: entrypointPath,
					Target: "seed:/usr/local/bin/docker-entrypoint.sh",
					Type: "bind",
				},
				{
					Source: nodekeyPath,
					Target: "seed:/tendermint/config/node_key.json",
					Type: "bind",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("created container!!!!!!!: ", seedContainer.Created)

	err = client.StartContainer(seedContainer.ID, seedContainer.HostConfig)
	if err != nil {
		return nil, err
	}
	fmt.Println("started container!!!!!!: ", seedContainer.ID)
	return seedContainer, nil
}
