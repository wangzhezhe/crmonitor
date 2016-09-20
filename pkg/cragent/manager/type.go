package manager

import (
	"github.com/crmonitor/pkg/cragent/manager/net"
)

//weird error for test when put the struct into the crtype/types.go
//unknown crtype.Ctncpu field 'SysTime' in struct literal

//the final value to be stored in influxdb
//Ctncpu

type Ctncpu struct {
	UserPercen float32 //the percentage (%) of user time in last interval time
	SysPercen  float32 //the percentage (%) of sys time in last interval time
	TotalTime  int
	SysTime    int
	UserTime   int
}

//Ctnmem
type Ctnmem struct {
	memUsage   float64 //rss+cache
	memUpLimit float64
	memPercen  float32
}

type Ctnblkio struct {
	currBlkrbytes int
	currBlkwbytes int
	rbytesps      int
	wbytesps      int
}

//Ctnnet
type Ctnnet struct {
	InterfacesMap map[string]net.InterfaceStats
}

//Ctnlabel

type CtnLabels struct {
	OriginalLabels  map[string]string
	CustomizeLabels map[string]string
	// customizelabel should include
	// hostip
	// containerid
	// hostport
}

type CtnBasic struct {
	FstPid      int
	Status      string
	ContaienrIP string
}

type ContainerResource struct {
	CtnBasic
	Ctncpu
	Ctnmem
	Ctnnet
	Ctnblkio
	CtnLabels
}

type ContaienrMetric struct {
	Cpu      string
	Memory   string
	Blkio    string
	Net      string
	MetaInfo string
}

//TODO It's better to put the metainfo into the etcd
type MetaInfo struct {
	ContainerID    string
	NodeName       string
	ContainerState string
	NodeState      string
}
