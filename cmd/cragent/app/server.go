package app

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"time"

	"golang.org/x/net/context"

	"strconv"

	"github.com/coreos/etcd/client"
	"github.com/crmonitor/pkg/api"
	"github.com/crmonitor/pkg/crtype"
	"github.com/crmonitor/pkg/register"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
	"github.com/crmonitor/pkg/util/event"
)

var (
	DefaultTTL        = time.Duration(time.Second * 60)
	DefaultHostip     string
	DefaultServerport = 9999
	//attention !
	//the agent should be started with the sudo privilege
	//the dockerclient could not access the unix:///var/run/docker.sock otherwise
	DefaultDockerdaemon = "unix:///var/run/docker.sock"
	Defaultrootkey      = "crmonitor"
	// etcd support the nest dir
)

type CRAgent struct {
	//the url should be like http://127.0.0.1:4001
	ETCD_URL   string
	Etcdclient client.Client
	TTL        time.Duration
	Hostip     string
}

func NewCRAgent() *CRAgent {

	agent := CRAgent{
		TTL: DefaultTTL,
	}

	return &agent
}

func Run(c *CRAgent) error {
	etcdclient, err := etcdclienttool.GetEtcdclient(c.ETCD_URL)
	c.Etcdclient = etcdclient
	if err != nil {
		log.Println("failed to register the crmagent")
	}

	//register the agent to the etcd
	Nodeinstance := c.Getregisternodeinstance()
	registerValue, err := json.Marshal(Nodeinstance)
	if err != nil {
		log.Println("error , failed to marshal Nodeinstance ," + err.Error())
	}
	registeroption := &register.Registeroption{
		TTL:        c.TTL,
		Etcdclient: c.Etcdclient,
		Key:        Defaultrootkey + "/" + "cragent",
		Value:      string(registerValue),
	}
	err = c.Doregister(registeroption)
	if err != nil {
		log.Println("error , could not register into etcd: ", err.Error())
	}

	//start the docker driver

	//collect and register the image info
	dockerclient, err := etcdclienttool.GetDockerClient(DefaultDockerdaemon)
	if err != nil {
		log.Println("error , failed yo get dockerclient", err)
	}
	register.Defaultrootkey = Defaultrootkey
	err = register.Imageregisterinit(register.Defaultrootkey, dockerclient, etcdclient)
	if err != nil {
		log.Println("error , failed to do the image init", err)
	}
	//collect and register the container info
	err = register.Containerregisterinit(register.Defaultrootkey, DefaultHostip, dockerclient, etcdclient)
	if err != nil {
		log.Println("error , failed to do the container init", err)
	}

	//create the event manager and start to listen the docker client
	eventmanager := event.Eventmanager{}
	//the go routine is ended with the main process
	go func() { eventmanager.Parsevent() }()

	//start the api server
	apiengine := api.Getengine()
	apiengine = api.Loadcragentapi(apiengine)
	apiengine.Run(":" + strconv.Itoa(DefaultServerport))
	return nil
}

func (c *CRAgent) AddFlags() error {
	flag.StringVar(&c.ETCD_URL, "etcd_url", "http://127.0.0.1:4001", "the url of etcd for crmonitor like http://127.0.0.1:4001")
	flag.StringVar(&DefaultHostip, "hostip", "", "the ip to register into the etcd")

	flag.Parse()
	log.Println("use the etcd_url,", c.ETCD_URL)
	if DefaultHostip == "" {
		return errors.New("error , the hostip could not be empty")
	}
	log.Printf("use the host ip %s to do the register\n", DefaultHostip)
	event.DefaultHostip = DefaultHostip
	event.Defaultdockerendpoint = DefaultDockerdaemon
	event.Defaultetcdurl = c.ETCD_URL
	return nil
}

func (c *CRAgent) Getregisternodeinstance() *crtype.Nodeinfo {
	Nodeinstance := &crtype.Nodeinfo{
		Hostip:       DefaultHostip,
		Agentport:    strconv.Itoa(DefaultServerport),
		Dockerdaemon: DefaultDockerdaemon,
	}

	return Nodeinstance
}

func (c *CRAgent) Doregister(option *register.Registeroption) error {

	etcdclient := option.Etcdclient
	if etcdclient == nil {
		return errors.New("error : faile to rigister to the etcd , etcd client could not be nil")
	}
	setoption := &client.SetOptions{
		TTL: c.TTL,
	}
	kapi := client.NewKeysAPI(etcdclient)
	if option.Key == "" || option.Value == "" {
		log.Printf("option %+v", option)
		return errors.New("error : faile to rigister to the etcd , key and value should not be empty")
	}

	//do a register option every TTL/2
	registerInterval := c.TTL / 2
	ticker := time.NewTicker(registerInterval)
	//this goroutine will be finished along with the main process
	go func() {
		for t := range ticker.C {
			_, err := kapi.Set(context.Background(), option.Key, option.Value, setoption)
			if err != nil {
				log.Println("error : faile to rigister to the etcd , " + err.Error())
				return
			}
			log.Println("succeeded in sending the register signal tick at : ", t)
		}
	}()

	return nil
}
