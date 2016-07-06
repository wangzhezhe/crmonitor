package manager

//weird error for test when put the struct into the crtype/types.go
//unknown crtype.Ctncpu field 'SysTime' in struct literal

//the final value to be stored in influxdb
//Ctncpu

type Ctncpu struct {
	UserPencen float32 //the percentage (%) of user time in last interval time
	SysPencen  float32 //the percentage (%) of sys time in last interval time
	TotalTime  int
	SysTime    int
	UserTime   int
}

//Ctnmem
type Ctnmem struct {
	memUsage int //rss+cache / total
}

//Ctnnet
type Ctnnet struct {
	rxBytesRate int
	txBytesRate int
}

//Ctnlabel

type CtnLabels struct {
	OriginalLabels  map[string]string
	CustomizeLabels map[string]string
}

type ContainerResource struct {
	Ctncpu
	Ctnmem
	Ctnnet
	CtnLabels
}
