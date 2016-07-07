package memory

import (
	"fmt"
	"io/ioutil"
	"math"

	"regexp"
	"strconv"
)

var (
	DefaultMemoInfoFile = "/proc/meminfo"

	MemSubCgroupPath       = "/memory/docker"
	memoryCapacityRegexp   = regexp.MustCompile(`MemTotal:\s*([0-9]+) kB`)
	swapCapacityRegexp     = regexp.MustCompile(`SwapTotal:\s*([0-9]+) kB`)
	cgroupMemIntegerRegexp = regexp.MustCompile(`([0-9]+)`)
)

type MemoryManager struct {
}

func CheckLimitation(limit float64) bool {
	//if the uplimit value is pow(2,64)-1
	//their is no limitation, return false
	highestValue := math.Pow(2, 64)
	if limit == highestValue {
		return false
	} else {
		return true
	}

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
	matches := cgroupMemIntegerRegexp.FindSubmatch(out)
	value, err := strconv.ParseInt(string(matches[0]), 10, 64)
	if err != nil {
		return -1, err
	}
	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(value), nil
}

//get memlimit for specific container
//if their is no limitation the value would be the 18446744073709551615
func (m *MemoryManager) GetContainerMemLimit(path string) (float64, error) {
	fileName := path + "/" + "memory.limit_in_bytes"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, err
	}
	matches := cgroupMemIntegerRegexp.FindSubmatch(out)
	value, err := strconv.ParseFloat(string(matches[0]), 64)

	if CheckLimitation(value) {
		return value, nil
	} else {
		fmt.Println("do not have the limitation, use /proc/meminfo")
		totalValue, err := m.GetTotalMem()
		if err != nil {
			return -1, err
		} else {
			return float64(totalValue), err
		}
	}
}

//this is the memtotal in /proc/meminfo
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
