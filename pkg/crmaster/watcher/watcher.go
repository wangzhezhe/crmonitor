package watcher

import (
	"encoding/json"
	"log"
	"strings"

	"net/http"

	"crypto/tls"
	"io/ioutil"

	"github.com/coreos/etcd/client"
	"github.com/crmonitor/pkg/crtype"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
	"golang.org/x/net/context"
)

var Defaultrootkey string

//&{Action:update Node:{Key: /crmonitor/images/postgres:latest/projectd/a2ae4c9535514c40d990ab8fe3785d5f054a8752f5b0a05aed09045bcdc8557c, CreatedIndex: 48820, ModifiedIndex: 48840, TTL: 0} PrevNode:{Key: /crmonitor/images/postgres:latest/projectd/a2ae4c9535514c40d990ab8fe3785d5f054a8752f5b0a05aed09045bcdc8557c, CreatedIndex: 48820, ModifiedIndex: 48820, TTL: 0} Index:48839}
func Parsekey(Action string, Event string) *crtype.Event {
	/*
		type Event struct {
			ProjectName string `json:"project_name"`
			ImageName   string `json:"image_name"`
			ContainerID string `json:"container_id"`
			Event       string `json:"event"`
		}
	*/

	eventdetail := &crtype.Event{}
	splitlist := strings.Split(Event, "/")
	for index, item := range splitlist {
		if item == "images" {
			eventdetail.Event = Action
			eventdetail.ImageName = splitlist[index+1]
			eventdetail.ProjectName = splitlist[index+2]
			eventdetail.ContainerID = splitlist[index+3]
		}
	}
	log.Println("the split list: ", splitlist)
	log.Printf("detail %+v ", eventdetail)
	return eventdetail

}

type Etcdwatcher struct {
}

//watch three kinds of operations
/*
2016/05/20 00:56:34 the watch info &{Action:update Node:{Key: /crmonitor/images/redis/tocontainers/bbcaa6d785ea7c10d0ff043db32fe2658a243ee630db4554dbcb125773351025, CreatedIndex: 47626, ModifiedIndex: 47754, TTL: 0} PrevNode:{Key: /crmonitor/images/redis/tocontainers/bbcaa6d785ea7c10d0ff043db32fe2658a243ee630db4554dbcb125773351025, CreatedIndex: 47626, ModifiedIndex: 47750, TTL: 0} Index:47750}
2016/05/20 00:56:34 prepare to listen
2016/05/20 00:56:34 the watch info &{Action:delete Node:{Key: /crmonitor/images/redis/tocontainers/bbcaa6d785ea7c10d0ff043db32fe2658a243ee630db4554dbcb125773351025, CreatedIndex: 47626, ModifiedIndex: 47755, TTL: 0} PrevNode:{Key: /crmonitor/images/redis/tocontainers/bbcaa6d785ea7c10d0ff043db32fe2658a243ee630db4554dbcb125773351025, CreatedIndex: 47626, ModifiedIndex: 47754, TTL: 0} Index:47754}
2016/05/20 00:56:34 prepare to listen
2016/05/20 00:56:47 the watch info &{Action:set Node:{Key: /crmonitor/images/erere, CreatedIndex: 47756, ModifiedIndex: 47756, TTL: 0} PrevNode:<nil> Index:47755}
*/
func (e *Etcdwatcher) StartEtcdwatcher(ETCD_URL string) {

	etcdclient, err := etcdclienttool.GetEtcdclient(ETCD_URL)
	if err != nil {
		log.Println("failed to get the client")
	}
	kapi := client.NewKeysAPI(etcdclient)
	watchkey := Defaultrootkey + "/" + "images"
	log.Println("the watchkey: ", watchkey)
	watcher := kapi.Watcher(watchkey, &client.WatcherOptions{Recursive: true})

	for {
		//log.Println("prepare to listen")
		respond, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("the err in watching: ", err)
		}
		log.Printf("the watch info %+v \n", respond)

		//parse the watch info
		//&{Action:update Node:{Key: /crmonitor/images/postgres:latest/projectd/a2ae4c9535514c40d990ab8fe3785d5f054a8752f5b0a05aed09045bcdc8557c, CreatedIndex: 48820, ModifiedIndex: 48840, TTL: 0} PrevNode:{Key: /crmonitor/images/postgres:latest/projectd/a2ae4c9535514c40d990ab8fe3785d5f054a8752f5b0a05aed09045bcdc8557c, CreatedIndex: 48820, ModifiedIndex: 48820, TTL: 0} Index:48839}
		watchkey := respond.Node.Key
		//log.Println("the ")
		eventinstance := Parsekey(respond.Action, watchkey)
		log.Println(eventinstance)
		//write the data to the channel
		//use post to replace that
		//104.236.190.90:8000/send
		url := "http://104.236.190.90:8000/send"

		jsonstr, _ := json.Marshal(eventinstance)
		reqest, err := http.NewRequest("POST", url, strings.NewReader(string(jsonstr)))
		if err != nil {
			panic(err)
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
		client := &http.Client{Transport: tr}
		response, err := client.Do(reqest)
		if err != nil {
			log.Println(err)
			return
		}
		returnbody, _ := ioutil.ReadAll(response.Body)

		log.Println(string(returnbody))

	}

}
