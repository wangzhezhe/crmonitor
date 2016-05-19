package clienttool

import (
	"github.com/fsouza/go-dockerclient"
)

//attention !!!
//the agent should be started with the sudo privilege
//the dockerclient could not access the unix:///var/run/docker.sock otherwise

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
