package app

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/crmonitor/cmd/cragent/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	c := app.NewCRAgent()
	err := c.AddFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = app.Run(c); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
