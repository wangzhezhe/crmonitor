package manager

import (
	"regexp"

	"github.com/crmonitor/pkg/cragent/manager/cpu"
	"github.com/crmonitor/pkg/cragent/manager/memory"
	"github.com/crmonitor/pkg/cragent/manager/net"

	"github.com/crmonitor/pkg/util/clienttool"
	"golang.org/x/net/context"
)

var (
	DefaultCgroupDir      = "/sys/fs/cgroup"
	DefaultDockerReg      = regexp.MustCompile(`[a-zA-Z0-9-_+.]+:[a-fA-F0-9]+`)
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

type ContaienrLabels map[string]string

type Manager struct {
	net.NetManager
	cpu.CpuManager
	memory.MemoryManager
}

//env is used to store the config info like mysql connection
//label is used to mark the identity of the container
//refer https://docs.docker.com/engine/userguide/labels-custom-metadata/ to check about the labels
func (m *Manager) GetLabels(containerID string) (ContaienrLabels, error) {
	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)
	if err != nil {
		return nil, err
	}
	inspectInfo, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}
	labels := inspectInfo.Config.Labels
	return labels, nil
}

var ContainerCurrnetDetail map[string]*ContainerResource

func init() {
	ContainerCurrnetDetail = make(map[string]*ContainerResource)
}

// get info at this point
func (m *Manager) GetContainerAllInfo(containerID string) (*ContainerResource, error) {
	lastContainerInfo, ok := ContainerCurrnetDetail[containerID]
	_ = lastContainerInfo
	var firstCheck bool
	firstCheck = false
	if ok == false {
		//first time to get the container info
		lastContainerInfo := &ContainerResource{}
		_ = lastContainerInfo
		firstCheck = true

	} else {
		lastContainerInfo = ContainerCurrnetDetail[containerID]
	}

	//cpu
	cpuInfopath := "/sys/fs/cgroup" + cpu.CpuacctSubCgroupPath + "/" + containerID
	currTotalTime, err := m.CpuManager.GetTotalCpuTime()
	if err != nil {
		return nil, err
	}
	currSysTime, currUserTime, err := m.CpuManager.GetCpuTimeInPath(cpuInfopath)
	if err != nil {
		return nil, err
	}
	currCpuInfo := Ctncpu{
		TotalTime: currTotalTime,
		SysTime:   currSysTime,
		UserTime:  currUserTime,
	}
	if firstCheck == false {
		sysPercen := 100 * float32(currSysTime-lastContainerInfo.Ctncpu.SysTime) / float32(currTotalTime-lastContainerInfo.Ctncpu.TotalTime)
		userPercen := 100 * float32(currSysTime-lastContainerInfo.Ctncpu.SysTime) / float32(currTotalTime-lastContainerInfo.Ctncpu.TotalTime)
		currCpuInfo.SysPencen = sysPercen
		currCpuInfo.UserPencen = userPercen
	}

	//mem

	//net

	//label

	//return currContainerInfo, nil

	currContainerInfo := &ContainerResource{
		Ctncpu: currCpuInfo,
	}

	//update the old container info
	ContainerCurrnetDetail[containerID] = currContainerInfo
	return currContainerInfo, nil

}

//start the ticker and send the info into the influx db
func (m *Manager) Start(contaienrList []string) {
	//merge the new container list into the old one
	//initiate at first time
	//then using notify to monitor
	//or get the container dir direactly

	//get the info

	//push into influx
}
