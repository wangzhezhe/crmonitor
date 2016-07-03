package clienttool

import (
	"fmt"
	"testing"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func TestGetDockerClient(t *testing.T) {
	/* deprecated for using old docker client
		t.Skip()
		endpoint := "unix:///var/run/docker.sock"
		dockerclient, err := GetDockerClient(endpoint)
		if err != nil {
			fmt.Println(err)
		}

		info, err := dockerclient.Info()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(info)
		//event example :
		//image option
		//{Status:pull ID:ubuntu:14.04 From: Time:1463214573}
		//{Status:untag ID:sha256:90d5884b1ee07f7f791f51bab92933943c87357bcd2fa6be0e82c48411bbb653 From: Time:1463214686}
		//{Status:delete ID:sha256:90d5884b1ee07f7f791f51bab92933943c87357bcd2fa6be0e82c48411bbb653 From: Time:1463214686}
		//container option
		//{Status:start ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214794}
		//{Status:die ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214860}
		//{Status:destroy ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214860}
		//get event , update the container status in path
		//use the full sha256 id to be the key of image
		eventchannel := make(chan *docker.APIEvents, 1)
		err = dockerclient.AddEventListener(eventchannel)
		if err != nil {
			fmt.Println(err)
		}
	A:
		for {
			select {
			case value, ok := <-eventchannel:
				if !ok {
					break A
				}
				fmt.Printf("get docker event %+v:", value)
			default:
				//fmt.Println("sleep 1s")
				time.Sleep(time.Second * 1)
			}

		}
	*/

}

func TestListContainer(t *testing.T) {

	endpoint := "unix:///var/run/docker.sock"
	cli, err := GetDockerClient(endpoint)
	if err != nil {
		t.Error(err)
	}
	// refer to https://docs.docker.com/v1.11/engine/reference/api/docker_remote_api/
	// to check the meaning of the option
	options := types.ContainerListOptions{All: true}
	containers, err := cli.client.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		fmt.Println(c.ID)
	}

}

//get the first pid in a container
func TestInspectContainer(t *testing.T) {
	endpoint := "unix:///var/run/docker.sock"
	cli, err := GetDockerClient(endpoint)
	if err != nil {
		t.Error(err)
	}
	containerID := "6d59e0d97082"
	inspectInfo, err := cli.client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("the inspectInfo: %+v\n", inspectInfo)
	fmt.Println("the Pid in container Stat: ", inspectInfo.State.Pid)
}
