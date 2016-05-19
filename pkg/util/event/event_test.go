package event

import (
	"testing"
	"time"

	"github.com/crmonitor/pkg/register"
)

/*

	listener := make(chan *APIEvents, 10)
	defer func() {
		time.Sleep(10 * time.Millisecond)
		if err := client.RemoveEventListener(listener); err != nil {
			t.Error(err)
		}
	}()

	err = client.AddEventListener(listener)
	if err != nil {
		t.Errorf("Failed to add event listener: %s", err)
	}

	timeout := time.After(1 * time.Second)
	var count int

	for {
		select {
		case msg := <-listener:
			t.Logf("Received: %v", *msg)
			count++
			err = checkEvent(count, msg)
			if err != nil {
				t.Fatalf("Check event failed: %s", err)
			}
			if count == 4 {
				return
			}
		case <-timeout:
			t.Fatalf("%s timed out waiting on events", testName)
		}
	}




*/

func TestParsevent(t *testing.T) {
	DefaultHostip = "10.211.55.14"
	Defaultdockerendpoint = "unix:///var/run/docker.sock"
	Defaultetcdurl = "http://127.0.0.1:4001"
	register.Defaultrootkey = "crmonitor"
	eventmanager := Eventmanager{}

	go func() { eventmanager.Parsevent() }()
	time.Sleep(time.Second * 100000)

}

/*
	dockerclient, err := clienttool.GetDockerClient(Defaultdockerendpoint)

	if err != nil {
		log.Println("failed to get dockerclient ", err)
	}
	//etcdclient, err := clienttool.GetEtcdclient(Defaultetcdurl)
	if err != nil {
		log.Println("failed to get etcdclient ", err)
	}
	eventchannel := make(chan *docker.APIEvents, 10)
	err = dockerclient.AddEventListener(eventchannel)
	if err != nil {
		log.Println("failed to add event listener : ", err)
	}
	log.Println("start the event listener")
	//A:
	for {
		//if get value from channel , ok is true
		//value, ok := <-eventchannel
		//log.Println("the value ", value)
		//if ok == true {
		select {

		case value, ok := <-eventchannel:
			if ok {
				log.Printf("get docker event %+v:", value)
			}

			//log.Println("ok?", ok)
			//if value.Status == "start" || value.Status == "die" || value.Status == "destroy" {
			//reister the new container status into etcd
			//rootkey string, status string, containerid string, repotag string
			//eventstatus string, containerid string, repotag string, hostip string, dockerclient *docker.Client
			//register.Containerinfoupdate(value.Status, value.ID, value.From, DefaultHostip, dockerclient, etcdclient)
			//} else {
			//	log.Println("other event info type")
			//}

		default:
			log.Println("sleep 1s")
			time.Sleep(time.Second * 1)

		}

	}*/
