package events

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	eventtypes "github.com/docker/engine-api/types/events"
)

func TestMonitor(t *testing.T) {
	//endpoint := "unix:///var/run/docker.sock"
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Error(err)
	}

	//ctx, cancel := context.WithCancel(context.Background())
	//_=cancel
	ctx := context.Background()
	errChan := Monitor(ctx, cli, types.EventsOptions{}, func(event eventtypes.Message) {
		fmt.Printf("get the docker event : %v\n", event)
	})

	if err := <-errChan; err != nil {
		fmt.Println("failed to get the event")
	}
}
