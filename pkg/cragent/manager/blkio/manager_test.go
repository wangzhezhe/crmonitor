package blkio

import (
	"log"
	"testing"
)

func TestGetBlkio(t *testing.T) {
	testcontaienrID := "8852499e82d44635081f1cc25cc6c735d6db7cce3875aa498d4702c9068a8b6a"
	DefaultCgroupDir := "/sys/fs/cgroup"
	path := DefaultCgroupDir + BlkioSubCgroupPath + "/" + testcontaienrID
	blkioManager := &BlkIOManager{}
	user, sys, err := blkioManager.GetBlkio(path)
	if err != nil {
		t.Error(err)
	}
	log.Println("readbytes: ", user, " write bytes: ", sys)

}
