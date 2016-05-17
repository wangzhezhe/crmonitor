package manager

import (
	"errors"
	"log"
	"strings"

	"encoding/json"

	"github.com/coreos/etcd/client"
	"github.com/crmonitor/pkg/crtype"
	"github.com/crmonitor/pkg/register"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
	"golang.org/x/net/context"
)

var (
	Defaultrootkey = "crmonitor"
	Defaultproject = "crproject"
	Defaultetcdurl string
)

type CRMastermanager struct {
	ETCD_URL   string
	Etcdclient client.Client
}

func GetCRMastermanager(etcd_url string) (*CRMastermanager, error) {
	//attention : add http !!!
	if !strings.Contains(etcd_url, "http") {
		etcd_url = "http://" + etcd_url
	}
	etcdclient, err := etcdclienttool.GetEtcdclient(etcd_url)
	if err != nil {
		return nil, err
	}
	log.Println("the url: ", Defaultetcdurl)
	if Defaultetcdurl == "" {
		return nil, errors.New("failed to get the addr of etcd")
	}
	return &CRMastermanager{
		ETCD_URL:   Defaultetcdurl,
		Etcdclient: etcdclient,
	}, nil
}

func getContainerlistfromimage(imagename string, etcdclient client.Client) ([]crtype.Container, error) {
	kapi := client.NewKeysAPI(etcdclient)
	key := Defaultrootkey + "/" + "images" + "/" + imagename + "/" + "tocontainers"
	rawdata, err := kapi.Get(context.Background(), key, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			log.Printf("error , do not have any containers record for %s \n ", imagename)
			return nil, nil
		} else {
			return nil, err
		}
	}

	containerlist := []crtype.Container{}
	for _, container := range rawdata.Node.Nodes {
		p := crtype.Container{}
		jsoninfo := container.Value
		//log.Println("the json info", jsoninfo)
		json.Unmarshal([]byte(jsoninfo), &p)
		//serche the container by image
		containerlist = append(containerlist, p)
	}

	return containerlist, nil
}

func (c *CRMastermanager) Getproject() (interface{}, error) {

	etcdclient := c.Etcdclient
	kapi := client.NewKeysAPI(etcdclient)
	key := Defaultrootkey + "/" + Defaultproject
	projectlist, err := kapi.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("the app list %+v", projectlist)
	plist := []crtype.CRMProject{}
	Nodevalue := projectlist.Node.Nodes
	for _, project := range Nodevalue {
		p := crtype.CRMProject{}
		jsoninfo := project.Value
		log.Println("the json info", jsoninfo)
		json.Unmarshal([]byte(jsoninfo), &p)
		//serche the container by image
		for index, layer := range p.Layers {
			containerlist, err := getContainerlistfromimage(layer.Imagename, etcdclient)
			if err != nil {
				log.Printf("error , failed to get container list for image : %s\n", layer.Imagename)
			}
			//log.Printf("the contianer list %+v", containerlist)
			p.Layers[index].Containerlist = containerlist
		}

		plist = append(plist, p)
	}

	return plist, nil
}

func (c *CRMastermanager) Registerproject(projectname string, rowdata string) error {

	registeroption := &register.Registeroption{
		Etcdclient: c.Etcdclient,
		Key:        Defaultrootkey + "/" + Defaultproject + "/" + projectname,
		Value:      rowdata,
	}
	err := c.Doregister(registeroption)
	if err != nil {
		return err
	}
	return nil
}

func (c *CRMastermanager) Doregister(option *register.Registeroption) error {
	etcdclient := option.Etcdclient
	if etcdclient == nil {
		return errors.New("error : failed to rigister to the etcd , etcd client could not be nil")
	}

	kapi := client.NewKeysAPI(etcdclient)
	if option.Key == "" || option.Value == "" {
		log.Printf("option %+v", option)
		return errors.New("error : faile to rigister to the etcd , key and value should not be empty")
	}
	_, err := kapi.Set(context.Background(), option.Key, option.Value, nil)
	if err != nil {
		log.Println("error : failed to rigister to the etcd for CRMastermanager , " + err.Error())
		return err
	}
	return nil
}
