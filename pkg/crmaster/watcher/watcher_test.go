package watcher

import (
	"testing"
)

func TestStartEtcdwatcher(t *testing.T) {
	testwatcher := &Etcdwatcher{}

	etcdurl := "http://127.0.0.1:4001"
	Defaultrootkey = "crmonitor"
	testwatcher.startEtcdwatcher(etcdurl)

}
