package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/rumblesan/distribnodes/pinger"
	"github.com/rumblesan/distribnodes/types"
)

func main() {
	fmt.Println("Distrib Nodes")

	cfg := types.AppConfig{
		PingDuration: time.Duration(30 * time.Second),
	}

	self := types.DistribNode{
		ID:      types.GenerateNodeID(),
		Address: fmt.Sprintf("localhost:%d", 3000),
	}

	as := types.AppState{
		Self:     self,
		AppMutex: &sync.Mutex{},
	}

	pinger.Start(&as, &cfg)

}
