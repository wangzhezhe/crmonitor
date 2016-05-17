package crtype

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetengine(t *testing.T) {

	var crmlist []CRMApp
	crmappa := CRMApp{
		Name: "testappa",
	}

	crmappb := CRMApp{
		Name: "testappb",
	}
	crmlist = append(crmlist, crmappa)
	crmlist = append(crmlist, crmappb)

	value, err := json.Marshal(crmlist)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value))

	var nodelist []Nodeinfo
	nodea := Nodeinfo{
		Hostip:       "192.168.0.1",
		Agentport:    "9999",
		Dockerdaemon: "192.168.0.1:6379",
	}
	nodeb := Nodeinfo{
		Hostip:       "192.168.0.2",
		Agentport:    "9999",
		Dockerdaemon: "192.168.0.2:6379",
	}
	nodelist = append(nodelist, nodea)
	nodelist = append(nodelist, nodeb)
	value, err = json.Marshal(nodelist)
	fmt.Println(string(value))
	var imagelist []Image
	imagea := Image{
		Repository: "repositorya",
		Tag:        "taga",
		Id:         "testid",
	}
	imageb := Image{
		Repository: "repositorya",
		Tag:        "taga",
		Id:         "testid",
	}
	imagelist = append(imagelist, imagea)
	imagelist = append(imagelist, imageb)
	value, err = json.Marshal(imagelist)
	fmt.Println(string(value))

}
