package register

import (
	"testing"

	"fmt"

	"github.com/crmonitor/pkg/util/clienttool"
)

func TestImageregister(t *testing.T) {

	//rootkey string, dockerclient *docker.Client, etcdclient *etcdclient.Client
	rootkey := "crmonitor"
	endpoint := "unix:///var/run/docker.sock"
	etcd_url := "http://127.0.0.1:4001"
	dockerclient, err := clienttool.GetDockerClient(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	etcdclient, err := clienttool.GetEtcdclient(etcd_url)
	if err != nil {
		fmt.Println(err)
	}
	err = Imageregisterinit(rootkey, dockerclient, etcdclient)
	if err != nil {
		fmt.Println(err)
	}
	hostip := "10.211.55.15"
	err = Containerregisterinit(rootkey, hostip, dockerclient, etcdclient)
	if err != nil {
		fmt.Println(err)
	}
}
