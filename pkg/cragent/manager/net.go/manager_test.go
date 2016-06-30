package memory

import (
	"log"
	"testing"
)

func TestGetNetInfo(t *testing.T) {
	netManager := &NetManager{}
	dev := "eth4"
	stat := "rx_bytes"
	value, err := netManager.GetNetInfo(dev, stat)
	if err != nil {
		t.Error(err)
	}
	log.Println("the value: ", value)

}
