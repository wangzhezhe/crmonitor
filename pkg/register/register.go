package register

import (
	"time"

	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/coreos/etcd/client"
	etcdclientpack "github.com/coreos/etcd/client"
	"github.com/fsouza/go-dockerclient"
	"golang.org/x/net/context"
)

type Registeroption struct {
	TTL        time.Duration
	Etcdclient client.Client
	Key        string
	Value      string
}

//the CRMAgent implement this interface
type Register interface {
	// Register to the etcd
	Doregister(option *Registeroption) error
}

var Imagerootpath = "images"
var Subimagedetailpath = "node"
var Subcontainerdetailpath = "tocontainers"

//get the image info on local machine and register into the etcd
//path: rootkey/image/<imagename>.../node/<imagedetail>
//                                  /tocontainers/<containername>/containerdetail
// the imageid is unique in different machine??
//do not set the ttl
//using watch mechanism to update the image list on etcd when pulling and deleting new images
func Imageregisterinit(rootkey string, dockerclient *docker.Client, etcdclient etcdclientpack.Client) error {
	imageinsertpath := rootkey + "/" + Imagerootpath
	imagelist, err := dockerclient.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return err
	}
	// merge the images in local machine into the etcd
	kapi := etcdclientpack.NewKeysAPI(etcdclient)
	for _, image := range imagelist {
		for _, repotag := range image.RepoTags {
			if strings.Contains(repotag, "none") {
				continue
			}
			//attention!
			//there are multiple dir layers if it is contains "/" in key value
			inseartpath_image := imageinsertpath + "/" + repotag + "/" + Subimagedetailpath

			resp, err := kapi.Get(context.Background(), inseartpath_image, nil)
			if resp != nil && err == nil {
				continue
			}
			log.Printf("image %s not exist in etcd , register it \n", repotag)
			jsonvalue, err := json.Marshal(image)
			if err != nil {
				return errors.New("error fail to do json marshal in Imageregisterinit : " + err.Error())
			}
			_, err = kapi.Set(context.Background(), inseartpath_image, string(jsonvalue), nil)

			if err != nil {
				return errors.New("error fail to do the image register : " + err.Error())
			}

		}
	}
	return nil

}

//get the container metadata
//register them into the rootkeyrootkey/image/includecontainer/imagename(repotag)/containermeta//imagename(repotag)/containermeta/
func Containerregisterinit(rootkey string, hostip string, dockerclient *docker.Client, etcdclient etcdclientpack.Client) error {
	imageinsertpath := rootkey + "/" + Imagerootpath
	//get the container on local machine
	containerlist, err := dockerclient.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		errors.New("failed to get container info in Containerregisterinit : " + err.Error())
	}

	kapi := etcdclientpack.NewKeysAPI(etcdclient)
	for _, container := range containerlist {
		repotag := container.Image

		inseartpath_container := imageinsertpath + "/" + repotag + "/" + Subcontainerdetailpath + "/" + container.ID
		// the original struct in docker client have been modified
		container.Hostip = hostip
		jsonvalue, err := json.Marshal(container)
		if err != nil {
			return errors.New("error , failed to do json marshal in Containerregisterinit : " + err.Error())
		}

		// option Dir could be used to assign wether the node is created as dir
		_, err = kapi.Set(context.Background(), inseartpath_container, string(jsonvalue), nil)

		if err != nil {
			return errors.New("error fail to do the image register : " + err.Error())
		}
	}

	return nil

}

// CURD according to the event message
// {Status:pull ID:ubuntu:14.04 From: Time:1463214573}
// {Status:untag ID:sha256:90d5884b1ee07f7f791f51bab92933943c87357bcd2fa6be0e82c48411bbb653 From: Time:1463214686}
// {Status:delete ID:sha256:90d5884b1ee07f7f791f51bab92933943c87357bcd2fa6be0e82c48411bbb653 From: Time:1463214686}
func Imageregisterevent(rootkey string, eventinfo docker.APIEvents, etcdclient etcdclientpack.Client) error {
	if eventinfo.Status == "pull" {
		//get image and merge
	} else if eventinfo.Status == "untag" {

	} else if eventinfo.Status == "delete" {
		//if there still exist container , do not delete
	} else {

	}

	return nil
}
