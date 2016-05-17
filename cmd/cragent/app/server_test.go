package app

import (
	"testing"

	//"github.com/coreos/etcd/client"
	"github.com/crmonitor/pkg/cragent/register"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
	//"golang.org/x/net/context"
	"encoding/json"
	"time"

	"github.com/crmonitor/pkg/api"

	//"golang.org/x/net/context"
)

/*
type CRAgent struct {
	ETCD_URL   string
	Etcdclient client.Client
	TTL        time.Duration
	Hosturl    string
}

type Registeroption struct {
	TTL        time.Duration
	Etcdclient client.Client
	Key        string
	Value      []byte
}
*/

func TestDoregister(t *testing.T) {

	//create the agent instance
	tmpETCD_URL := "http://127.0.0.1:4001"
	tmpTTL := time.Duration(time.Second * 60)
	tmpClient, err := etcdclienttool.GetEtcdclient(tmpETCD_URL)
	if err != nil {
		t.Error("fail to create etcdclient")
	}
	cragent := NewCRAgent()
	cragent.ETCD_URL = tmpETCD_URL
	cragent.Etcdclient = tmpClient
	cragent.TTL = tmpTTL
	cragent.Hosturl = "127.0.0.1"
	cragent.Etcdclient = tmpClient

	//create the setoption
	tmpMap := map[string]string{"testa": "valuea", "testb": "valueb"}
	tmpValue, err := json.Marshal(tmpMap)
	if err != nil {
		t.Error("fail to marshal tmpMap")
	}
	registeroption := &register.Registeroption{
		TTL:        tmpTTL,
		Etcdclient: tmpClient,
		Key:        "testregister/subdira/1234567",
		Value:      string(tmpValue),
	}

	//do the register operation
	err = cragent.Doregister(registeroption)
	if err != nil {
		t.Error(err.Error())
	}
	/*
		//check the return value
		kapi := client.NewKeysAPI(tmpClient)
		resp, err := kapi.Get(context.Background(), registeroption.Key, nil)
		if err != nil {
			t.Error(err.Error())
		}

		t.Logf("respond node value %+v", resp.Node.Value)
	*/

	time.Sleep(time.Second * 100)
}

func TestGetregisternodeinstance(t *testing.T) {
	t.Skip()
	Nodeinstance := &api.Nodeinfo{
		Hostip:       "127.0.0.1",
		Agentport:    "9999",
		Dockerdaemon: "test.sock",
	}
	t.Logf("the node value:%+v", Nodeinstance)
	value, err := json.Marshal(Nodeinstance)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(value))
}
