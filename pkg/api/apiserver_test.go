package api

import (
	"fmt"
	"testing"
)

func TestGetCRMastermanager(t *testing.T) {
	manager, err := getCRMastermanager()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("manager %+v", manager)
	}
}
