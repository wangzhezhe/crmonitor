package clienttool

import (
	"time"

	"github.com/coreos/etcd/client"
)

var Etcdclient client.Client

func GetEtcdclient(ETCD_URL string) (client.Client, error) {
	if Etcdclient != nil {
		return Etcdclient, nil
	}
	cfg := client.Config{
		Endpoints: []string{ETCD_URL},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}
