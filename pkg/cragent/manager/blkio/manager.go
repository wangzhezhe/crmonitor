package blkio

import (
	"io/ioutil"
	//"log"
	"regexp"
	"strconv"
)

var (
	BlkioSubCgroupPath = "/blkio/docker"
	BlkioReadRegexp    = regexp.MustCompile(`8:0 Read\s*([0-9]+)`)
	BlkioWriteRegexp   = regexp.MustCompile(`8:0 Write\s*([0-9]+)`)
)

type BlkIOManager struct {
}

//the major:minor 8:0
//8:0 Read 26968064
//8:0 Write 0

//return read bytes write bytes err
func (b *BlkIOManager) GetBlkio(path string) (int, int, error) {
	fileName := path + "/" + "blkio.throttle.io_service_bytes"
	out, err := ioutil.ReadFile(fileName)
	if err != nil {
		return -1, -1, err
	}
	matches := BlkioReadRegexp.FindSubmatch(out)
	//debug
	//log.Println("blkio matches: ", matches)
	//log.Println("file name: ", fileName)
	var readBytes int
	if len(matches) == 0 {
		readBytes = 0
	} else {
		value, err := strconv.ParseInt(string(matches[1]), 10, 64)
		if err != nil {
			return -1, -1, err
		}
		readBytes = int(value)
	}

	matches = BlkioWriteRegexp.FindSubmatch(out)
	var writeBytes int
	if len(matches) == 0 {
		writeBytes = 0
	} else {
		value, err := strconv.ParseInt(string(matches[1]), 10, 64)
		if err != nil {
			return -1, -1, err
		}
		writeBytes = int(value)
	}

	//log.Printf("the mem capacity for %s is %d", path, value)
	return int(readBytes), int(writeBytes), nil

}
