package pinger

import (
	"log"
	"time"

	"github.com/rumblesan/distribnodes/types"
)

// Start will start the pinger
func Start(as *types.AppState, cfg *types.AppConfig) *time.Ticker {
	t := time.NewTicker(cfg.PingDuration)

	go Pinger(as, t)

	return t
}

// Pinger handles sending Ping messages to remote nodes
func Pinger(as *types.AppState, ticker *time.Ticker) {
	for t := range ticker.C {
		log.Printf("%s - Pinging Nodes\n", t)

		nodeList := as.GetNodeList()
		for _, node := range nodeList {
			err := pingNode(node)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func pingNode(n types.DistribNode) error {
	log.Printf("Pinging %s\n", n)
	return nil
}
