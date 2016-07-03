package clienttool

import (
	dockerclient "github.com/docker/engine-api/client"
	//"github.com/docker/engine-api/types"
	//"golang.org/x/net/context"
)

//attention !!! use engine api to interact with docker
//the agent should be started with the sudo privilege
//the dockerclient could not access the unix:///var/run/docker.sock otherwise

type DockerClient struct {
	client *dockerclient.Client
}

var DefaultDockerClient *DockerClient

func GetDockerClient(endpoint string) (*DockerClient, error) {
	if DefaultDockerClient != nil {
		return DefaultDockerClient, nil
	} else {
		defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
		cli, err := dockerclient.NewClient(endpoint, "v1.22", nil, defaultHeaders)
		if err != nil {
			return nil, err
		}

		newClient := &DockerClient{
			client: cli,
		}
		DefaultDockerClient = newClient

	}

	return DefaultDockerClient, nil

}
