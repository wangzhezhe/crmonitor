package manager

import (
	"errors"
	"log"
	"regexp"
	"strconv"

	"time"

	"github.com/crmonitor/cmd/cragent/conf"
	"github.com/crmonitor/pkg/cragent/manager/blkio"
	"github.com/crmonitor/pkg/cragent/manager/cpu"
	"github.com/crmonitor/pkg/cragent/manager/memory"
	"github.com/crmonitor/pkg/cragent/manager/net"
	"github.com/crmonitor/pkg/util/clienttool"
	"github.com/crmonitor/pkg/util/events"
	"github.com/docker/docker/api/types"
	eventtypes "github.com/docker/docker/api/types/events"
	"golang.org/x/net/context"
)

/*
	env stores the connection info, like the password or the addresses of other services
	in k8s it will be like

	labels store the identity info like env=production and project.name=abc


				"Labels": {
	                "io.kubernetes.container.hash": "6059dfa2",
	                "io.kubernetes.container.name": "POD",
	                "io.kubernetes.container.restartCount": "0",
	                "io.kubernetes.container.terminationMessagePath": "",
	                "io.kubernetes.pod.name": "podtest",
	                "io.kubernetes.pod.namespace": "default",
	                "io.kubernetes.pod.terminationGracePeriod": "30",
	                "io.kubernetes.pod.uid": "533b3f11-4439-11e6-a14d-001c422a82c8"
	            }
				 "Env": [
                "KUBERNETES_SERVICE_HOST=10.0.0.1",
                "KUBERNETES_SERVICE_PORT=443",
                "KUBERNETES_SERVICE_PORT_HTTPS=443",
                "KUBERNETES_PORT=tcp://10.0.0.1:443",
                "KUBERNETES_PORT_443_TCP=tcp://10.0.0.1:443",
                "KUBERNETES_PORT_443_TCP_PROTO=tcp",
                "KUBERNETES_PORT_443_TCP_PORT=443",
                "KUBERNETES_PORT_443_TCP_ADDR=10.0.0.1"
            ],

			in docker-compose way it would be like:
			            "Labels": {
                "com.docker.compose.config-hash": "b43b53fefbeb61f5154d9eb9b1b7e2a38598c9cfbb8c7ea88b5d56a3b19dff1d",
                "com.docker.compose.container-number": "1",
                "com.docker.compose.oneoff": "False",
                "com.docker.compose.project": "poc",
                "com.docker.compose.service": "mysql",
                "com.docker.compose.version": "1.7.0"
            }

			           "Env": [
                "MYSQL_ROOT_PASSWORD=password",
                "MYSQL_PASSWORD=password",
                "MYSQL_USER=sstack",
                "MYSQL_DATABASE=sstack",
                "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
                "MYSQL_MAJOR=5.6",
                "MYSQL_VERSION=5.6.28-1debian8"
            ],

how to assign a label for the docker container ? refer to
https://docs.docker.com/engine/userguide/labels-custom-metadata/
use container label to support mutiple orchesatratin service

*/

var (
	DefaultCgroupDir      = "/sys/fs/cgroup"
	DefaultDockerReg      = regexp.MustCompile(`[a-zA-Z0-9-_+.]+:[a-fA-F0-9]+`)
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

type Manager struct {
	net.NetManager
	cpu.CpuManager
	memory.MemoryManager
	blkio.BlkIOManager
	Interval int
}

var (
	ContainerCurrnetDetail map[string]*ContainerResource
	ContainerMetricName    ContaienrMetric
)

func init() {
	//how to add lock for ContainerCurrnetDetail ???
	ContainerCurrnetDetail = make(map[string]*ContainerResource)
	ContainerMetricName = ContaienrMetric{
		Cpu:    "cpu",
		Memory: "memory",
		Net:    "net",
		Blkio:  "blkio",
	}
}

func (m *Manager) GetCustomizeLabel(containerID string, inspectInfo types.ContainerJSON) map[string]string {
	labels := make(map[string]string)
	labels["containerID"] = containerID
	labels["hostip"] = conf.GlobalConfig.DefaultHostip
	//labels["status"] = inspectInfo.ContainerJSONBase.State.Status

	return labels

}

//env is used to store the config info like mysql connection
//label is used to mark the identity of the container
//refer https://docs.docker.com/engine/userguide/labels-custom-metadata/ to check about the labels

func (m *Manager) GetLabels(containerID string) (map[string]string, map[string]string, error) {
	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)
	if err != nil {
		return nil, nil, err
	}
	inspectInfo, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, nil, err
	}
	containerlabels := inspectInfo.Config.Labels

	customizelabels := make(map[string]string)

	customizelabels = m.GetCustomizeLabel(containerID, inspectInfo)

	customizelabels["image"] = inspectInfo.Image

	return containerlabels, customizelabels, nil
}

//register all the container info into ContainerCurrnetDetail (add labels)
func (m *Manager) InitialContainer() error {
	//get all the container id
	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)
	if err != nil {
		return err
	}
	options := types.ContainerListOptions{All: true}
	containerList, err := dockerClient.ContainerList(context.Background(), options)
	if err != nil {
		return err
	}

	//inspect each container to get the labels, and stored the label info into the map
	for _, containerItem := range containerList {
		containerID := containerItem.ID
		labels := containerItem.Labels
		//debug
		log.Println("the labels: ", labels)
		_, ok := ContainerCurrnetDetail[containerID]

		if ok == false {
			originalLabels, customizeLabels, err := m.GetLabels(containerID)
			if err != nil {
				return err
			}
			tmpDetail := &ContainerResource{CtnLabels: CtnLabels{OriginalLabels: originalLabels, CustomizeLabels: customizeLabels}}
			ContainerCurrnetDetail[containerID] = tmpDetail
		}
	}
	return nil

}

var ADD = "add"
var DELETE = "delete"

func (m *Manager) AdjustContainer(containerID string, action string) error {
	//if the action is add, add the container into ContainerCurrnetDetail
	if action == ADD {
		_, ok := ContainerCurrnetDetail[containerID]
		if ok == false {
			originalLabels, customizeLabels, err := m.GetLabels(containerID)
			if err != nil {
				return err
			}
			//if the action is restart the container may already exist

			tmpDetail := &ContainerResource{CtnLabels: CtnLabels{OriginalLabels: originalLabels, CustomizeLabels: customizeLabels}}

			ContainerCurrnetDetail[containerID] = tmpDetail
			log.Println("add the container with ID " + containerID)
		}

	} else if action == DELETE {
		//if the action is delete, delete the contaienr from ContainerCurrnetDetail
		delete(ContainerCurrnetDetail, containerID)
		log.Println("delete the container with ID " + containerID)
	} else {
		return errors.New("invalid action: " + action)
	}
	return nil
}

//refer to the https://docs.docker.com/engine/reference/api/docker_remote_api/#docker-events
//for all event changes
//only consider the add and destroyed
func (m *Manager) parseEvent(event eventtypes.Message) {
	if event.Status == "start" {
		//if the key with the containerID does not exist in the ContainerCurrnetDetails, add it
		containerID := event.ID
		m.AdjustContainer(containerID, "add")
	} else if event.Status == "destroy" {
		//remove the key with containerID form ContainerCurrnetDetails
		containerID := event.ID
		m.AdjustContainer(containerID, "delete")

	} else {
		//neglect events for other status change
		return
	}

	//debug
	//log.Printf("the current map: %+v\n", ContainerCurrnetDetail)
	return

}

//detect the event of docker daemon
//send
func (m *Manager) ParseEventInfo() error {

	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)

	if err != nil {
		return errors.New("failed to create docker client")
	}
	ctx := context.Background()
	errChan := events.Monitor(ctx, dockerClient, types.EventsOptions{}, m.parseEvent)
	if err := <-errChan; err != nil {
		log.Println("failed to get the event")
	}
	return err
}

func (m *Manager) HousKeeping(interval int, containerDetailMap map[string]*ContainerResource) {
	log.Println("start houskeeping")
	//range the container map in every interval second
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	for t := range ticker.C {
		log.Println("Tick at", t)
		//range the map and update the info
		for key, _ := range containerDetailMap {
			newInfo, err := m.GetContainerAllInfo(key, interval)
			if err != nil {
				//if the container is not started it will be failed to get info
				//only the running container have the record in cgroups
				log.Println("failed to get container info for: ", key, err.Error())
				continue
			}
			//transfer info into influxdata
			influxDataList, err := m.TransferIntoInfluxData(key, newInfo)
			if err != nil {
				log.Println("failed to transfer container data into influxdata for error: ", key, err.Error())
				continue
			}
			//send influxdata into influx db
			//log.Printf("new container info for %s: %+v\n ", key, newInfo)
			//TODO get info from env
			dbName := conf.GlobalConfig.DefaultInfluxDBContainer
			influxClient, err := clienttool.GetinfluxClient(dbName)
			if err != nil {
				log.Println("failed to get the influxclient")
				continue
			}
			err = influxClient.AddStats(influxDataList)
			if err != nil {
				log.Println("failed to upload influx data to influxdb: ", key, err.Error())
				continue
			}
		}
	}

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

	//basic info
	dockerClient, err := clienttool.GetDockerClient(DefaultDockerEndpoint)
	if err != nil {
		return nil, err
	}
	currFirstPid, currStatus, err := dockerClient.GetBasicInfoFromContainerID(containerID)
	if err != nil {
		return nil, err
	}
	currBasicInfo := CtnBasic{
		FstPid: currFirstPid,
		Status: currStatus,
	}

	if currStatus != "running" {

		currContainerInfo := &ContainerResource{

			CtnLabels: lastContainerInfo.CtnLabels, //use the last label info
			CtnBasic:  currBasicInfo,
		}

		//update the old container info
		ContainerCurrnetDetail[containerID] = currContainerInfo
		return currContainerInfo, nil

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
		currCpuInfo.SysPercen = sysPercen
		currCpuInfo.UserPercen = userPercen
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

	//blkio
	blkioPath := DefaultCgroupDir + blkio.BlkioSubCgroupPath + "/" + containerID

	currBlkRBytes, currBlkWBytes, err := m.BlkIOManager.GetBlkio(blkioPath)
	if err != nil {
		return nil, err
	}

	rbytesps := (currBlkRBytes - lastContainerInfo.Ctnblkio.currBlkrbytes) / interval
	wbytesps := (currBlkWBytes - lastContainerInfo.Ctnblkio.currBlkwbytes) / interval
	currBlkIOInfo := Ctnblkio{
		currBlkrbytes: currBlkRBytes,
		currBlkwbytes: currBlkWBytes,
		rbytesps:      rbytesps,
		wbytesps:      wbytesps,
	}
	//debug
	log.Printf("the blkioinfo %+v : ", currBlkIOInfo)

	//net
	//currFstPid := lastContainerInfo.CtnBasic.FstPid

	lastInterfaceMap := lastContainerInfo.Ctnnet.InterfacesMap
	currInterfaceMap, err := m.NetManager.GetNetInfoFromProc(currFirstPid)
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
		Ctnblkio:  currBlkIOInfo,
		CtnLabels: lastContainerInfo.CtnLabels, //use the last label info
		CtnBasic:  currBasicInfo,
	}

	//update the old container info
	ContainerCurrnetDetail[containerID] = currContainerInfo
	return currContainerInfo, nil

}

//start the ticker and send the info into the influx db
func (m *Manager) Start() error {

	//initiate at first time
	err := m.InitialContainer()
	if err != nil {
		log.Println("error for Start option: ", err)
		return err
	}

	//go func start parsing the log info
	go func() {
		err := m.ParseEventInfo()
		if err != nil {
			log.Println("error for Parsing event: ", err)
			return
		}
	}()

	//stat the ticker and get the container detail info
	m.HousKeeping(m.Interval, ContainerCurrnetDetail)

	//push container info into influxdb
	return nil
}

//every ContainerResource could be transfered into multiple influxData
func (m *Manager) TransferIntoInfluxData(containerID string, resource *ContainerResource) ([]*clienttool.InfluxData, error) {

	var influxList []*clienttool.InfluxData
	tags := resource.OriginalLabels
	for k, v := range resource.CustomizeLabels {
		tags[k] = v
	}
	//add the container status and first pid
	tags["firstpid"] = strconv.Itoa(resource.CtnBasic.FstPid)

	tags["status"] = resource.CtnBasic.Status

	//cpu
	measurement := ContainerMetricName.Cpu
	fields := map[string]interface{}{
		"syspercen":  resource.Ctncpu.SysPercen,
		"userpercen": resource.Ctncpu.UserPercen,
		"totaltime":  resource.Ctncpu.TotalTime,
	}

	influxData := &clienttool.InfluxData{Measurement: measurement, Fields: fields, Tags: tags}
	influxList = append(influxList, influxData)

	//mem

	measurement = ContainerMetricName.Memory
	fields = map[string]interface{}{
		"memuplimit": resource.Ctnmem.memUpLimit,
		"memusage":   resource.Ctnmem.memUsage,
		"mempercern": resource.Ctnmem.memPercen,
	}

	influxData = &clienttool.InfluxData{Measurement: measurement, Fields: fields, Tags: tags}
	influxList = append(influxList, influxData)

	//blkio
	measurement = ContainerMetricName.Blkio
	fields = map[string]interface{}{
		"blkiorps": resource.Ctnblkio.rbytesps,
		"blkiowps": resource.Ctnblkio.wbytesps,
	}

	influxData = &clienttool.InfluxData{Measurement: measurement, Fields: fields, Tags: tags}
	influxList = append(influxList, influxData)

	//net
	for key, value := range resource.Ctnnet.InterfacesMap {

		measurement = containerID + "_" + key + "_" + ContainerMetricName.Net
		fields = map[string]interface{}{
			"rxbytes":     value.RxBytes,
			"txbytes":     value.TxBytes,
			"rxbytesrate": value.RxBytesRate,
			"txbytesrate": value.TxBytesRate,
		}

		influxData = &clienttool.InfluxData{Measurement: measurement, Fields: fields, Tags: tags}
		influxList = append(influxList, influxData)
	}

	return influxList, nil

}
