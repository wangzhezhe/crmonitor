package memory

import (
	"fmt"
	"io/ioutil"

	"regexp"
	"strconv"
)

var (
	DefaultMemoInfoFile     = "/proc/meminfo"
	MemSubCgroupPath        = "/memory/docker"
	memoryCapacityRegexp    = regexp.MustCompile(`MemTotal:\s*([0-9]+) kB`)
	swapCapacityRegexp      = regexp.MustCompile(`SwapTotal:\s*([0-9]+) kB`)
	cgroupMemCapacityRegexp = regexp.MustCompile(`([0-9]+)`)
)

type MemoryManager struct {
}

// copy from cadvisor
// parseCapacity matches a Regexp in a []byte, returning the resulting value in bytes.
// Assumes that the value matched by the Regexp is in KB.
func parseCapacity(b []byte, r *regexp.Regexp) (int, error) {
	matches := r.FindSubmatch(b)
	if len(matches) != 2 {
		return -1, fmt.Errorf("failed to match regexp in output: %q", string(b))
	}
	m, err := strconv.ParseInt(string(matches[1]), 10, 64)
	if err != nil {
		return -1, err
	}

	// Convert to bytes.
	return int(m * 1024), err
}

//path should be like /sys/fs/cgroup/memory/docker/<containerid>
func (m *MemoryManager) GetContainerMemCapacity(path string) (int, error) {
	fileName := path + "/" + "memory.usage_in_bytes"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, err
	}
	matches := cgroupMemCapacityRegexp.FindSubmatch(out)
	value, err := strconv.ParseInt(string(matches[0]), 10, 64)
	if err != nil {
		return -1, err
	}
	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(value), nil
}

func (m *MemoryManager) GetTotalMem() (int, error) {

	out, err := ioutil.ReadFile(DefaultMemoInfoFile)
	if err != nil {
		return -1, err
	}

	memoryCapacity, err := parseCapacity(out, memoryCapacityRegexp)
	if err != nil {
		return -1, err
	}

	return memoryCapacity, nil

}
