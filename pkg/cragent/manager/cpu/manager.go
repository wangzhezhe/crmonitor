package cpu

import (
	"io/ioutil"

	"regexp"
	"strconv"
)

var (
	CpuacctSubCgroupPath = "/cpuacct/docker" //get cpuacct.stat from here (nm)
	cpuUserRegexp        = regexp.MustCompile(`user\s*([0-9]+)`)
	cpuSysRegexp         = regexp.MustCompile(`system\s*([0-9]+)`)
)

type CpuManager struct {
}

//return usertime systemtime error
func (c *CpuManager) GetCpuTimeInPath(path string) (int, int, error) {
	fileName := path + "/" + "cpuacct.stat"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, -1, err
	}
	matches := cpuUserRegexp.FindSubmatch(out)
	userTime, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	matches = cpuSysRegexp.FindSubmatch(out)
	sysTime, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(userTime), int(sysTime), nil
}

//return the total usage
func (c *CpuManager) GetTotalCpuTime() (int, error) {
	fileName := "/sys/fs/cgroup/cpuacct/cpuacct.stat"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, err
	}
	matches := cpuUserRegexp.FindSubmatch(out)
	userTime, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, err
	}
	matches = cpuSysRegexp.FindSubmatch(out)
	sysTime, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, err
	}
	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(userTime) + int(sysTime), nil
}
