package packet

import (
	"time"
)

type Packetdetail struct {
	Requestdetail string
	Responddetail string
}
type HttpTransaction struct {
	//insert struct
	Packetdetail
	Srcip       string
	Srcport     string
	Destip      string
	Destport    string
	Timesend    time.Time
	Timereceive time.Time
	Respondtime float64
	//only application layer info
}

type Address struct {
	IP   string
	PORT string
}
