package clienttool

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func TestGetEtcdclient(t *testing.T) {
	t.Skip()
	etcd_url := "http://127.0.0.1:4001"
	etcdclient, err := GetEtcdclient(etcd_url)
	if err != nil {
		t.Error("failed to get the client")
	}
	kapi := client.NewKeysAPI(etcdclient)
	resp, err := kapi.Get(context.Background(), "/abcde", nil)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("respond %+v", resp)
	t.Logf("respond %+V", resp)
	option := &client.SetOptions{
		TTL: time.Duration(time.Second * 60),
	}

	resp, err = kapi.Set(context.Background(), "/test60s", "hello", option)
	if err != nil {
		t.Error(err.Error())
	}
	//time.Sleep(time.Second * 30)
	resp, err = kapi.Set(context.Background(), "/test60s", "world ", option)
	t.Logf("the respond value \n %+v", resp)
	resp, err = kapi.Get(context.Background(), "/testetcd", nil)
	if err != nil {
		t.Error(err.Error())
	}

	if resp.Node.Value != "testvalue" {
		t.Error(errors.New("return value failure"))
	}
	resp, err = kapi.Set(context.Background(), "/abc'/'def", "world ", nil)

}
func TestWatcher(t *testing.T) {
	etcd_url := "http://127.0.0.1:4001"
	etcdclient, err := GetEtcdclient(etcd_url)
	if err != nil {
		t.Error("failed to get the client")
	}
	kapi := client.NewKeysAPI(etcdclient)
	watchkey := "/crmonitor/images"
	watcher := kapi.Watcher(watchkey, &client.WatcherOptions{Recursive: true})

	respond, err := watcher.Next(context.Background())
	if err != nil {
		fmt.Println("the err: ", err)
	}
	fmt.Printf("the watch info %+v \n", respond)

}
