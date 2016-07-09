package conf

import (
	"os"
	"strconv"
)

//get all the config info from env or the config file
type Config struct {
	Interval              int
	DefaultHostip         string
	DefaultDockerEndpoint string
	DefaultEtcdURL        string
	DefaultInfluxURL      string
	DefaultInfluxDB       string
}

var GlobalConfig *Config

func init() {
	intervalStr := os.Getenv("INTERVALTIME")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		panic("the env INTERVAL should be an integer")
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
	defaultInfluxDB := os.Getenv("INFLUXDBNAME")
	if defaultInfluxDB == "" {
		panic("the env INFLUXDBNAME should be empty")
	}

	GlobalConfig = &Config{
		Interval:              interval,
		DefaultHostip:         hostIP,
		DefaultDockerEndpoint: DockerEndPoint,
		DefaultEtcdURL:        DefaultEtcdUrl,
		DefaultInfluxURL:      defaultInfluxURL,
		DefaultInfluxDB:       defaultInfluxDB,
	}
}
