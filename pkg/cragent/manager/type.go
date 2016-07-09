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
	FstPid int
}

type ContainerResource struct {
	CtnBasic
	Ctncpu
	Ctnmem
	Ctnnet
	CtnLabels
}

type ContaienrMetric struct {
	Cpu    string
	Memory string
	Net    string
}
