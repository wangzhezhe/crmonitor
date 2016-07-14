package packet

import (
	"testing"
	"time"
)

func TestStartCollect(t *testing.T) {
	t.Log("test packet collection")
	//testInterface := "eth5"
	testInterface := "lo"
	interval := 2 * time.Second
	expiration := time.After(time.Second * 600)
	StartCollect(testInterface, interval, expiration)

}
