package conf

import (
	"os"
	"strconv"
)

//get all the config info from env or the config file
type Config struct {
	Interval int
}

var GlobalConfig *Config

func init() {
	intervalStr := os.Getenv("INTERVALTIME")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		panic("the env INTERVAL should be an integer")
	}

	GlobalConfig = &Config{
		Interval: interval,
	}
}
