package packet

import (
	"log"

	"net"
	"strings"

	"container/list"
	"sync"
	"time"

	"github.com/crmonitor/cmd/cragent/conf"
	"github.com/crmonitor/pkg/util/clienttool"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

//install libpcap before using sudo apt-get install libpcap0.8-dev
var (
	snapshotLen int32 = 65535
	promiscuous bool  = false

	handle             *pcap.Handle
	localip            string
	httpinstancelist   *list.List
	Activeflag         bool = false //if flag is true , the agent is collecting the data
	Flagmutex               = &sync.Mutex{}
	DefaultInterface   string
	globalInfluxClient *clienttool.InfluxdbStorage
)

func init() {
	dbname := conf.GlobalConfig.DefaultInfluxDBPacket
	client, err := clienttool.GetinfluxClient(dbname)
	globalInfluxClient = client
	if err != nil {
		log.Println("failed to init the influx client")
	}
}

func GetLocalip(iface string) (string, error) {
	ifaceobj, err := net.InterfaceByName(iface)
	if err != nil {
		return "", err
	}
	addrarry, err := ifaceobj.Addrs()
	if err != nil {
		return "", err
	}
	var localip = ""

	log.Println(addrarry)

	for _, ip := range addrarry {
		IP := ip.String()
		// attention to the lo interface
		if iface == "lo" {
			localip = strings.TrimSuffix(IP, "/8")
			return localip, nil
		}
		if strings.Contains(IP, "/24") {
			localip = strings.TrimSuffix(IP, "/24")
		}
	}

	return localip, nil
}

//detect the http packet return the info
func detectHttp(packet gopacket.Packet) (bool, []byte) {
	applicationLayer := packet.ApplicationLayer()
	//panic will happenden if try to transfer into string
	//log.Println("the content of application layer", string(applicationLayer.Payload()))
	if applicationLayer != nil {
		// do not support to the packet using https
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {

			log.Println("HTTP found!")

			return true, applicationLayer.LayerContents()
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

//if it is the output stream from local machine
func outputStream(packet gopacket.Packet, Srcaddr *Address, Destaddr *Address) {

	ishttp, httpcontent := detectHttp(packet)
	_ = httpcontent

	if ishttp {
		sendtime := time.Now()
		//iphandler := packet.Layer(layers.LayerTypeIPv4)
		reqdetail := string(packet.ApplicationLayer().LayerContents())
		httpinstance := &HttpTransaction{
			Srcip:        Srcaddr.IP,
			Srcport:      Srcaddr.PORT,
			Destip:       Destaddr.IP,
			Destport:     Destaddr.PORT,
			Timesend:     sendtime,
			Packetdetail: Packetdetail{Requestdetail: reqdetail, Responddetail: ""},
		}
		//put the httpinstance into a list

		log.Printf("store the instance:%+v\n", httpinstance)

		httpinstancelist.PushBack(httpinstance)

		log.Printf("the length of the list : %+v\n", httpinstancelist.Len())

	}

}

//adjust if this is the response of the packet
func ifreverse(httpinstance *HttpTransaction, Srcaddr *Address, Destaddr *Address) bool {
	if httpinstance.Srcip == Destaddr.IP && httpinstance.Destip == Srcaddr.IP {
		if httpinstance.Srcport == Destaddr.PORT && httpinstance.Destport == Srcaddr.PORT {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//if it is the input stream to local machine
func inputStream(packet gopacket.Packet, Srcaddr *Address, Destaddr *Address) {
	//get the instance from the list which has the reverse srcaddr and the destaddr
	respdetail := string(packet.Data())

	log.Println("the length of the list before extract element:", httpinstancelist.Len())

	for element := httpinstancelist.Front(); element != nil; element = element.Next() {
		httpinstance := element.Value.(*HttpTransaction)
		isreverse := ifreverse(httpinstance, Srcaddr, Destaddr)
		if isreverse {
			httpinstance.Timereceive = time.Now()
			//Attention the units of respondtime is ms !!!
			responsetime := httpinstance.Timereceive.Sub(httpinstance.Timesend).Seconds() * 1000
			httpinstance.Respondtime = responsetime
			//store the respond detail
			//log.Println("respond info:", respdetail)
			httpinstance.Packetdetail.Responddetail = respdetail

			log.Printf("Respond duration:%vms\n", responsetime)
			log.Printf("Get the response: %v\n", httpinstance)

			httpinstancelist.Remove(element)
			//??how to use generic to realize the push function of different type
			//jsoninfo, _ := json.Marshal(httpinstance)
			//the type should be the ip of this machine
			//the first parameter is index the second one is type
			//err := ESClient.Push(jsoninfo, "packetagent", localip)
			//create the InfluxData
			//data := &clienttool.InfluxData{
			//	Measurement: "respondtime",
			//}
			log.Println("send the respondtime measurement")

			influxData := getInfluxDataFromHttpInstance(httpinstance)
			err := globalInfluxClient.AddStats(influxData)
			if err != nil {
				log.Println("failed to upload data into influxdb: ", err)
			}

			break
		}
	}

}

//return measurement,fields,tags
func getInfluxDataFromHttpInstance(httpinstance *HttpTransaction) []*clienttool.InfluxData {
	measurement := "respondtime"

	fields := make(map[string]interface{})
	fields["respondtime"] = httpinstance.Respondtime

	tags := make(map[string]string)
	tags["destip"] = httpinstance.Destip
	tags["destport"] = httpinstance.Destport
	tags["srcip"] = httpinstance.Srcip
	tags["srcport"] = httpinstance.Srcport
	data := &clienttool.InfluxData{
		Measurement: measurement,
		Fields:      fields,
		Tags:        tags,
	}

	var dataBatch []*clienttool.InfluxData

	dataBatch = append(dataBatch, data)
	return dataBatch
}

// get a new packet every time
func processPacketInfo(packet gopacket.Packet, localip string) {
	//get the specified layer
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {

		log.Println("TCP layer is detected.")

		tcphandler, _ := tcpLayer.(*layers.TCP)
		srcport := tcphandler.SrcPort
		destport := tcphandler.DstPort
		//get the specified layer
		iplayer := packet.Layer(layers.LayerTypeIPv4)
		httphandler, _ := iplayer.(*layers.IPv4)
		srcip := httphandler.SrcIP
		destip := httphandler.DstIP
		//log.Println(srcip.String())
		//send the packet from local machine
		Srcaddr := &Address{IP: srcip.String(), PORT: srcport.String()}
		Destaddr := &Address{IP: destip.String(), PORT: destport.String()}

		log.Printf("srcaddr %v destaddr %v \n", Srcaddr, Destaddr)

		var mutex = &sync.Mutex{}
		log.Println("the srcip.string", srcip.String(), "the localip", localip)
		if srcip.String() == localip {
			//mutex.Lock()
			outputStream(packet, Srcaddr, Destaddr)
			//mutex.Unlock()
		}
		//get the packet from the local machine
		if destip.String() == localip {

			mutex.Lock()
			inputStream(packet, Srcaddr, Destaddr)
			mutex.Unlock()
		}

	}
}

//device is the name of interface
func StartCollect(device string, interval time.Duration, timesignal <-chan time.Time) {
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, interval)
	if err != nil {
		log.Println(err.Error())
	}

	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	localIP, err := GetLocalip(device)
	if err != nil {
		log.Println("failed to get localip", localIP)
	}
	log.Println("the local IP:", localIP)

	httpinstancelist = list.New()
	if err != nil {

		log.Println(err.Error())
	}
A:
	for packet := range packetSource.Packets() {
		select {
		case <-timesignal:
			//stop the falg
			Flagmutex.Lock()
			Activeflag = false
			Flagmutex.Unlock()
			break A
		default:
			processPacketInfo(packet, localIP)
		}

	}
}
