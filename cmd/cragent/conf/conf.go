package conf

import (
	"log"
	"os"
	"strconv"
)

//get all the config info from env or the config file
type Config struct {
	Interval                 int
	DefaultHostip            string
	DefaultDockerEndpoint    string
	DefaultEtcdURL           string
	DefaultInfluxURL         string
	DefaultInfluxDBContainer string
	DefaultInfluxDBPacket    string
}

var GlobalConfig *Config

func init() {
	intervalStr := os.Getenv("INTERVALTIME")
	if intervalStr == "" {
		panic("the env INTERVALTIME should not be empty")
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Fatal(err)
	}

	if !(interval < 10 && interval > 0) {
		panic("the env INTERVALTIME should between 0-9")
	}

	hostIP := os.Getenv("HOSTIP")
	if hostIP == "" {
		hostIP = "127.0.0.1"
	}
	DockerEndPoint := os.Getenv("DOCKERENDPOINT")
	if DockerEndPoint == "" {
		DockerEndPoint = "unix:///var/run/docker.sock"
	}
	DefaultEtcdUrl := os.Getenv("ETCDURL")

	defaultInfluxURL := os.Getenv("INFLUXDBURL")
	if defaultInfluxURL == "" {
		defaultInfluxURL = "http://127.0.0.1:8086"
	}
	defaultInfluxDBContainer := os.Getenv("INFLUXDBNAME_CONTAINER")
	if defaultInfluxDBContainer == "" {
		panic("the env INFLUXDBNAME_CONTAINER should be empty")
	}
	defaultInfluxDBPacket := os.Getenv("INFLUXDBNAME_PACKET")
	if defaultInfluxDBPacket == "" {
		panic("the env INFLUXDBNAME_PACKET should be empty")
	}

	GlobalConfig = &Config{
		Interval:                 interval,
		DefaultHostip:            hostIP,
		DefaultDockerEndpoint:    DockerEndPoint,
		DefaultEtcdURL:           DefaultEtcdUrl,
		DefaultInfluxURL:         defaultInfluxURL,
		DefaultInfluxDBContainer: defaultInfluxDBContainer,
		DefaultInfluxDBPacket:    defaultInfluxDBPacket,
	}
}
