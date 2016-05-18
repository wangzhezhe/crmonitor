package event

import (
	"testing"
	"time"
)

func TestParsevent(t *testing.T) {
	DefaultHostip = "10.211.55.14"
	Defaultdockerendpoint = "unix:///var/run/docker.sock"
	Defaultetcdurl = "http://127.0.0.1:4001"
	eventmanager := Eventmanager{}
	//eventmanager.Parsevent()
	go func() { eventmanager.Parsevent() }()
	time.Sleep(time.Second * 100000)
}
