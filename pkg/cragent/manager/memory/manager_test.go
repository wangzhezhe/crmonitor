package memory

import (
	"log"
	"testing"
)

func TestGetContainerMemCapacity(t *testing.T) {
	testcontaienrID := "096eeebb39f13b271c9fdae8a8221a0b92f1c4a21e946d4a9d64a73b7c72aa50"
	DefaultCgroupDir := "/sys/fs/cgroup"
	path := DefaultCgroupDir + MemSubCgroupPath + "/" + testcontaienrID
	memManager := &MemoryManager{}
	value, err := memManager.GetContainerMemCapacity(path)
	if err != nil {
		t.Error(err)
	}
	log.Println("rss+cache: ", value)

	//check the mem limitation
	memUplimit, err := memManager.GetContainerMemLimit(path)
	if err != nil {
		t.Error(err)
	}
	log.Println("the up limitation: ", memUplimit)
}
