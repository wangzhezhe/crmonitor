package app

import (
	"flag"
	"log"
	"time"

	"strconv"

	"github.com/coreos/etcd/client"
	"github.com/crmonitor/pkg/api"
	"github.com/crmonitor/pkg/crmaster/manager"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
)

var (
	DefaultTTL          = time.Duration(time.Second * 60)
	DefaultHostip       string
	DefaultServerport   = 9998
	DefaultDockerdaemon = "unix:///var/run/docker.sock"
	Defaultrootkey      = "crmonitor"
	// etcd support the nest dir
)

type CRMaster struct {
	//the url should be like http://127.0.0.1:4001
	ETCD_URL   string
	Etcdclient client.Client
	TTL        time.Duration
	Hosturl    string
}

func NewCRMaster() *CRMaster {

	master := CRMaster{
		TTL: DefaultTTL,
	}

	return &master
}

func Run(c *CRMaster) error {
	etcdclient, err := etcdclienttool.GetEtcdclient(c.ETCD_URL)
	c.Etcdclient = etcdclient
	if err != nil {
		log.Println("failed to register the crmagent")
	}

	//register the agent to the etcd

	//start the docker driver

	//collect the image info

	//register the image info into the etcd

	//start the api server
	apiengine := api.Getengine()
	apiengine = api.Loadcrmasterapi(apiengine)
	apiengine.Run(":" + strconv.Itoa(DefaultServerport))
	return nil
}

func (c *CRMaster) AddFlags() error {
	flag.StringVar(&c.ETCD_URL, "etcd_url", "127.0.0.1:4001", "the url of etcd for crmonitor")
	//crmaster do not need to register itself into etcd
	flag.Parse()
	log.Println("use the etcd_url,", c.ETCD_URL)
	manager.Defaultetcdurl = c.ETCD_URL
	return nil
}
