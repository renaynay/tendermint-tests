package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

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
	time.Sleep(time.Minute*2)
	fmt.Println("killing container....")
	err = client.StopContainer(seed.ID, 20)
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
		},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("created container!!!!!!!: ", seedContainer.Created)

	// upload to container

	if err := uploadToContainer(client, seedContainer.ID, "docker-entrypoint.sh", entrypointPath, "/usr/local/bin/"); err != nil {
		return nil, err
	}
	if err := uploadToContainer(client, seedContainer.ID, "node_key.json", nodekeyPath,  "/tendermint/config/"); err != nil {
		return nil, err
	}

	err = client.StartContainer(seedContainer.ID, seedContainer.HostConfig)
	if err != nil {
		return nil, err
	}
	fmt.Println("started container!!!!!!: ", seedContainer.ID)
	return seedContainer, nil
}

// TODO make this more abstract
func uploadToContainer(client *docker.Client, containerID, filename, filepath, destination string) error {
	// Create a tarball archive with all the data files
	tarball := new(bytes.Buffer)
	tw := tar.NewWriter(tarball)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// Insert the file into the tarball archive
	header := &tar.Header{
		Name: filename,
		Mode: int64(0777),
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tw.Write(data); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	// Upload the tarball into the destination container
	return client.UploadToContainer(containerID, docker.UploadToContainerOptions{
		Context:     ctx,
		InputStream: tarball,
		Path:        destination,
	})
}
