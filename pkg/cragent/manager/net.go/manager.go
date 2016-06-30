package memory

import (
	"fmt"
	"io/ioutil"
	"path"
)

var (
	CpuacctSubCgroupPath = "/cpuacct/docker" //get cpuacct.stat from here (nm)
	DefaultNetDir        = "/sys/class/net"
)

type NetManager struct {
}

//stat could be the specific info for the device dev
//Interfaces     []InterfaceStats -> specific info
//containerid->interfaces?
func (*NetManager) GetNetInfo(dev string, stat string) (int, error) {
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

//get the veth of specific container
//refer to https://bugzilla.redhat.com/show_bug.cgi?id=1251538
//http://unix.stackexchange.com/questions/224201/what-is-proc-pid-net-dev
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

*/
