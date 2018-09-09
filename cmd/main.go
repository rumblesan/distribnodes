package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/rumblesan/distribnodes/config"
	"github.com/rumblesan/distribnodes/dnodes"
	"github.com/rumblesan/distribnodes/pinger"
	"github.com/rumblesan/distribnodes/types"
)

func main() {
	log.Println("Distrib Nodes")

	cfg := config.Get()

	self := types.DistribNode{
		ID:      dnodes.GenerateNodeID(),
		Address: fmt.Sprintf(":%d", cfg.NodePort),
	}

	log.Printf("Running as node %s\n", self.ID)

	app := &dnodes.DNodesApp{
		Self:     &self,
		NodeList: make(map[string]*types.NodeState),
		AppMutex: &sync.Mutex{},
	}

	for _, addr := range cfg.InitialNodes {
		err := pinger.FirstPing(app, addr)
		if err != nil {
			log.Println(err)
		}
	}

	pinger.Start(app, cfg)

	err := dnodes.StartRPCServer(app, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
