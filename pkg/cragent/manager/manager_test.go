package manager

import (
	"fmt"
	"testing"
	"time"
)

func TestGetLabels(t *testing.T) {
	manager := &Manager{}
	containerID := "6d59e0d9"
	labels, err := manager.GetLabels(containerID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("the labels %+v\n: ", labels)

}

func TestGetContainerAllInfo(t *testing.T) {
	manager := &Manager{}
	containerID := "096eeebb39f13b271c9fdae8a8221a0b92f1c4a21e946d4a9d64a73b7c72aa50"
	containerInfo, err := manager.GetContainerAllInfo(containerID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("container info: %+v ", containerInfo)
	time.Sleep(time.Second * 10)
	containerInfo, err = manager.GetContainerAllInfo(containerID)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("container info after interval second: %+v ", containerInfo)

}
