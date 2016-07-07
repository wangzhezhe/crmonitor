package manager

import (
	"log"
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

//register all the container info into ContainerCurrnetDetail (add labels)
func (m *Manager) RegisterContainer() {

}

// get info at this point
func (m *Manager) GetContainerAllInfo(containerID string, interval int) (*ContainerResource, error) {
	lastContainerInfo, ok := ContainerCurrnetDetail[containerID]
	var firstCheck bool
	firstCheck = false

	if ok == false {
		//first time to get the container info
		lastContainerInfo = &ContainerResource{}
		firstCheck = true

	} else {
		lastContainerInfo = ContainerCurrnetDetail[containerID]
	}

	//cpu
	cpuInfopath := DefaultCgroupDir + cpu.CpuacctSubCgroupPath + "/" + containerID
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
		userPercen := 100 * float32(currUserTime-lastContainerInfo.Ctncpu.UserTime) / float32(currTotalTime-lastContainerInfo.Ctncpu.TotalTime)
		currCpuInfo.SysPencen = sysPercen
		currCpuInfo.UserPencen = userPercen
	}

	//mem
	memInfoPath := DefaultCgroupDir + memory.MemSubCgroupPath + "/" + containerID

	currMemUsage, err := m.MemoryManager.GetContainerMemCapacity(memInfoPath)
	if err != nil {
		return nil, err
	}
	memUpLimit, err := m.MemoryManager.GetContainerMemLimit(memInfoPath)
	log.Println("the total mem: ", memUpLimit)
	if err != nil {
		return nil, err
	}
	memPercen := 100 * float32(currMemUsage) / float32(memUpLimit)

	currMemInfo := Ctnmem{
		memUsage:   float64(currMemUsage),
		memUpLimit: memUpLimit,
		memPercen:  memPercen,
	}

	//net
	//currFstPid := lastContainerInfo.CtnBasic.FstPid
	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)
	if err != nil {
		return nil, err
	}
	currFstPid, err := dockerClient.GetPidFromContainerID(containerID)
	if err != nil {
		return nil, err
	}
	lastInterfaceMap := lastContainerInfo.Ctnnet.InterfacesMap
	currInterfaceMap, err := m.NetManager.GetNetInfoFromProc(currFstPid)
	if err != nil {
		return nil, err
	}
	tmpMap := currInterfaceMap
	for key, value := range tmpMap {
		if firstCheck == false {

			rxBytesRate := float32(value.RxBytes-lastInterfaceMap[key].RxBytes) / float32(interval)
			txBytesRate := float32(value.TxBytes-lastInterfaceMap[key].TxBytes) / float32(interval)
			//update the value in map
			tmpInterface := currInterfaceMap[key]
			//could not assign the struct in map directly
			tmpInterface.RxBytesRate = rxBytesRate
			tmpInterface.TxBytesRate = txBytesRate
			currInterfaceMap[key] = tmpInterface

		}
	}

	currNetInfo := Ctnnet{
		InterfacesMap: currInterfaceMap,
	}

	//label (the label info should be detected and added into the container info when register the contaienr at the first time)

	//return currContainerInfo, nil

	currContainerInfo := &ContainerResource{
		Ctncpu:    currCpuInfo,
		Ctnmem:    currMemInfo,
		Ctnnet:    currNetInfo,
		CtnLabels: lastContainerInfo.CtnLabels, //use the last label info
		CtnBasic:  lastContainerInfo.CtnBasic,
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
