package event

import (
	"log"
	"time"

	"github.com/crmonitor/pkg/register"
	"github.com/crmonitor/pkg/util/clienttool"
	"github.com/fsouza/go-dockerclient"
)

var DefaultHostip string
var Defaultdockerendpoint string
var Defaultetcdurl string

type Eventmanager struct {
}

//{Status:start ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214794}
//{Status:die ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214860}
//{Status:destroy ID:538ef372f9d2a2402d48fe61c1cffd3c5519db78b3b198203d5068fcd4deeb70 From:busybox:latest Time:1463214860}
// get event , update the container status in path
func (e *Eventmanager) Parsevent() {
	dockerclient, err := clienttool.GetDockerClient(Defaultdockerendpoint)

	if err != nil {
		log.Println("failed to get dockerclient ", err)
	}
	etcdclient, err := clienttool.GetEtcdclient(Defaultetcdurl)
	if err != nil {
		log.Println("failed to get etcdclient ", err)
	}
	eventchannel := make(chan *docker.APIEvents, 5)
	err = dockerclient.AddEventListener(eventchannel)
	if err != nil {
		log.Println("failed to add event listener : ", err)
	}
	log.Println("start the event listener")
	//A:
	for {
		//if get value from channel , ok is true
		value, ok := <-eventchannel
		log.Println("the value ", value)
		if ok == true {

			log.Printf("get docker event %+v:", value)
			if value.Status == "start" || value.Status == "die" || value.Status == "destroy" || value.Status == "create" {
				//reister the new container status into etcd
				//rootkey string, status string, containerid string, repotag string
				//eventstatus string, containerid string, repotag string, hostip string, dockerclient *docker.Client
				log.Println("refreash the container")
				//quit with the main process
				register.Containerinfoupdate(value.Status, value.ID, value.From, DefaultHostip, dockerclient, etcdclient)
			} else {
				log.Println("other event info type")
			}

		} else {
			log.Println("sleep 1s")
			time.Sleep(time.Second * 1)
		}

	}

}
