package lib

import (
	"log"
	"testing"
)

func TestGetInfofromPortbylsof(t *testing.T) {
	port := 8099
	value, err := GetInfofromPortbylsof(port)
	if err != nil {
		log.Println(err)
	}
	log.Println("the value:", value)
}

func TestParseCgroupFile(t *testing.T) {
	port := 3322
	ifcontrol, value, err := IfinCgroupControl(port)
	if err != nil {
		log.Println(err)
	}
	log.Println(ifcontrol, value)

	port = 29769
	ifcontrol, value, err = IfinCgroupControl(port)
	if err != nil {
		log.Println(err)
	}
	log.Println(ifcontrol, value)

	//if answer is true , get container id
	//if answer is false, case a :the process is docker-proxy , and get the connected container ip from map
	//                    case b :other common process

}
