package watcher

import (
	"log"

	"github.com/coreos/etcd/client"
	etcdclienttool "github.com/crmonitor/pkg/util/clienttool"
	"golang.org/x/net/context"
)

var Defaultrootkey string

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
func (e *Etcdwatcher) startEtcdwatcher(ETCD_URL string) {

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

		//write the data to the channel

	}

}
