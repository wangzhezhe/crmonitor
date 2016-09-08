package net

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	CpuacctSubCgroupPath = "/cpuacct/docker" //get cpuacct.stat from here (nm)
	DefaultNetDir        = "/sys/class/net"
	DefaultProc          = "/proc"

	ignoredDevicePrefixes = []string{"lo", "veth", "docker", "br"} //do not show info with these  prefix
)

type InterfaceStats struct {
	// The name of the interface.
	Name string `json:"name"`
	// Cumulative count of bytes received.
	RxBytes uint64 `json:"rx_bytes"`
	// Cumulative count of packets received.
	RxPackets uint64 `json:"rx_packets"`
	// Cumulative count of receive errors encountered.
	RxErrors uint64 `json:"rx_errors"`
	// Cumulative count of packets dropped while receiving.
	RxDropped uint64 `json:"rx_dropped"`
	// Cumulative count of bytes transmitted.
	TxBytes uint64 `json:"tx_bytes"`
	// Cumulative count of packets transmitted.
	TxPackets uint64 `json:"tx_packets"`
	// Cumulative count of transmit errors encountered.
	TxErrors uint64 `json:"tx_errors"`
	// Cumulative count of packets dropped while transmitting.
	TxDropped   uint64  `json:"tx_dropped"`
	RxBytesRate float32 `json:"rx_bytes_rate"`
	TxBytesRate float32 `json:"tx_bytes_rate"`
}

type NetManager struct {
}

// stat could be the specific info for the device dev
// Interfaces     []InterfaceStats -> specific info
// containerid->interfaces?
func (*NetManager) GetNetInfoFromSys(dev string, stat string) (int, error) {
	statPath := path.Join(DefaultNetDir, dev, "/statistics", stat)

	out, err := ioutil.ReadFile(statPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read stat from %q for device %q", statPath, dev)
	}
	var s int
	n, err := fmt.Sscanf(string(out), "%d", &s)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("could not parse value from %q for file %s", string(out), statPath)
	}
	return s, nil
}

// input the pid and get the array of interfacesstats in this namespace
// refer to libcontainers/helpers.go
// show network info in same namespace
func (*NetManager) GetNetInfoFromProc(pid int) (map[string]InterfaceStats, error) {

	netStatsFile := path.Join(DefaultProc, strconv.Itoa(pid), "/net/dev")

	ifaceStats, err := scanInterfaceStats(netStatsFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read network stats: %v", err)
	}
	interfacesMap := make(map[string]InterfaceStats)
	for _, item := range ifaceStats {
		interfacesMap[item.Name] = item
	}

	return interfacesMap, nil

}
func isIgnoredDevice(ifName string) bool {
	for _, prefix := range ignoredDevicePrefixes {
		if strings.HasPrefix(strings.ToLower(ifName), prefix) {
			return true
		}
	}
	return false
}
func setInterfaceStatValues(fields []string, pointers []*uint64) error {
	for i, v := range fields {
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		*pointers[i] = val
	}
	return nil
}

func scanInterfaceStats(netStatsFile string) ([]InterfaceStats, error) {
	file, err := os.Open(netStatsFile)
	if err != nil {
		return nil, fmt.Errorf("failure opening %s: %v", netStatsFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Discard header lines
	for i := 0; i < 2; i++ {
		if b := scanner.Scan(); !b {
			return nil, scanner.Err()
		}
	}

	stats := []InterfaceStats{}
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, ":", "", -1)

		fields := strings.Fields(line)
		// If the format of the  line is invalid then don't trust any of the stats
		// in this file.
		if len(fields) != 17 {
			return nil, fmt.Errorf("invalid interface stats line: %v", line)
		}

		devName := fields[0]
		if isIgnoredDevice(devName) {
			continue
		}

		i := InterfaceStats{
			Name: devName,
		}

		statFields := append(fields[1:5], fields[9:13]...)
		statPointers := []*uint64{
			&i.RxBytes, &i.RxPackets, &i.RxErrors, &i.RxDropped,
			&i.TxBytes, &i.TxPackets, &i.TxErrors, &i.TxDropped,
		}

		err := setInterfaceStatValues(statFields, statPointers)
		if err != nil {
			return nil, fmt.Errorf("cannot parse interface stats (%v): %v", err, line)
		}

		stats = append(stats, i)
	}

	return stats, nil
}

// get the veth of specific container
// refer to https://bugzilla.redhat.com/show_bug.cgi?id=1251538
// http://unix.stackexchange.com/questions/224201/what-is-proc-pid-net-dev
// /proc/pid/net/dev shows the interface in same net namespace (from point of process)
// refer to libcontainer to get the net work info

/*
same data actually
root@ubuntu:/var/run/docker/netns# cat /proc/5408/net/dev
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0:   11715     108    0    0    0     0          0         0      648       8    0    0    0     0       0          0
    lo:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
root@ubuntu:/var/run/docker/netns# cat /sys/class/net/veth3154216/statistics/tx__bytes
cat: /sys/class/net/veth3154216/statistics/tx__bytes: No such file or directory
root@ubuntu:/var/run/docker/netns# cat /sys/class/net/veth3154216/statistics/tx_
tx_aborted_errors    tx_carrier_errors    tx_dropped           tx_fifo_errors       tx_packets
tx_bytes             tx_compressed        tx_errors            tx_heartbeat_errors  tx_window_errors
root@ubuntu:/var/run/docker/netns# cat /sys/class/net/veth3154216/statistics/tx_bytes
11715
root@ubuntu:/var/run/docker/netns#


refer docker-proxy
http://windsock.io/tag/docker-proxy/

listenports-> docker-proxy pid

gap:
when inspect the container , we could get the IPAddress of the contaienr
when we inspect the docker-proxy we could get the dest ip of this proxy
then we could make some index operaion to get the container id
containerid-> firstpid

userland-proxy to controle the proxy start

see the disadvantage of docker-proxy:http://www.dataguru.cn/thread-544489-1-1.html

input: the port
output: the containerid and the hostip

case a: using net==host, find pid from lsof -> /proc/cgroup(if in cgroup control) -> get contaienr id

case b: using default mode listen pid -> docker-proxy -> find dest ip -> find containerid by index map

*/
