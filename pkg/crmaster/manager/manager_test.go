package manager

import (
	"fmt"
	"testing"

	"github.com/crmonitor/pkg/util/clienttool"
)

//crmonitor/images/daocloud.io/daocloud/sstack_mysql:latest
//crmonitor/images/daocloud.io/daocloud/sstack_mysql:latest/tocontainers
func TestGetContainerlistfromimage(t *testing.T) {
	etcdclient, err := clienttool.GetEtcdclient("http://127.0.0.1:4001")
	if err != nil {
		fmt.Println(err)
	}
	imagename := "daocloud.io/daocloud/sstack_mysql:latest"
	containerlist, err := getContainerlistfromimage(imagename, etcdclient)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the container list %+v \n", containerlist)

}
