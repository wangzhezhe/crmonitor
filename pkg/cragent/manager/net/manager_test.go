package net

import (
	"log"
	"testing"
)

func TestGetNetInfo(t *testing.T) {
	netManager := &NetManager{}
	dev := "eth4"
	stat := "rx_bytes"
	value, err := netManager.GetNetInfoFromSys(dev, stat)
	if err != nil {
		t.Error(err)
	}
	log.Println("the value: ", value)

}

func TestGetNetInfoFromProc(t *testing.T) {
	//container with the default docker network
	pida := 15014
	netManager := &NetManager{}
	interfaceArray, err := netManager.GetNetInfoFromProc(pida)
	if err != nil {
		t.Error(err)
	}
	log.Printf("the interface array for pid %d is: %+v \n", pida, interfaceArray)

	//container without the network=host
	pidb := 15014
	interfaceArray, err = netManager.GetNetInfoFromProc(pidb)
	if err != nil {
		t.Error(err)
	}
	log.Printf("the interface array for pid %d is: %+v \n", pidb, interfaceArray)

}
