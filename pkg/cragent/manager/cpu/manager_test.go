package cpu

import (
	"log"
	"testing"
)

func TestGetContainerMemCapacity(t *testing.T) {
	testcontaienrID := "096eeebb39f13b271c9fdae8a8221a0b92f1c4a21e946d4a9d64a73b7c72aa50"
	DefaultCgroupDir := "/sys/fs/cgroup"
	path := DefaultCgroupDir + CpuacctSubCgroupPath + "/" + testcontaienrID
	cpuManager := &CpuManager{}
	user, sys, err := cpuManager.GetCpuTimeInPath(path)
	if err != nil {
		t.Error(err)
	}
	log.Println("usertime: ", user, " systime: ", sys)

	totalUsage, err := cpuManager.GetTotalCpuTime()
	if err != nil {
		t.Error(err)
	}
	log.Println("total usagetile(user+sys): ", totalUsage)

}
