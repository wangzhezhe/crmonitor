package clienttool

import (
	"github.com/fsouza/go-dockerclient"
)

var Dockerclient *docker.Client

func GetDockerClient(endpoint string) (*docker.Client, error) {
	if Dockerclient != nil {
		return Dockerclient, nil
	}
	Dockerclient, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return Dockerclient, nil
}
