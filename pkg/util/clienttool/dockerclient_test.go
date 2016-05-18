package clienttool

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func TestGetDockerClient(t *testing.T) {
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

}

func TestGetImage(t *testing.T) {
	//t.Skip()
	endpoint := "unix:///var/run/docker.sock"
	dockerclient, err := GetDockerClient(endpoint)
	image, err := dockerclient.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		fmt.Println(err)
	}
	value, _ := json.Marshal(image)
	fmt.Printf("list %+v", string(value))
	fmt.Println("+++++++++++++++++++")
	//imagedetail, err := dockerclient.InspectImage("94df4f")
	containerdetail, err := dockerclient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fmt.Println(err)
	}
	value, _ = json.Marshal(containerdetail)
	fmt.Printf("\n container deatils %+v \n", string(value))
}
