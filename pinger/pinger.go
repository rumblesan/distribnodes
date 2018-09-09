package pinger

import (
	"log"
	"time"

	"github.com/rumblesan/distribnodes/config"
	"github.com/rumblesan/distribnodes/dnodes"
	"github.com/rumblesan/distribnodes/types"
)

// Start will start the pinger
func Start(app *dnodes.DNodesApp, cfg *config.AppConfig) *time.Ticker {
	t := time.NewTicker(cfg.PingDuration)

	go Pinger(app, t)

	return t
}

// Pinger handles sending Ping messages to remote nodes
func Pinger(app *dnodes.DNodesApp, ticker *time.Ticker) {
	for t := range ticker.C {
		log.Printf("%s - Pinging Nodes\n", t)

		nodeList := app.GetRemoteNodeList()
		for _, node := range nodeList {
			err := PingNode(app, node)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// PingNode will send a ping message to a node and handle the response
func PingNode(app *dnodes.DNodesApp, n *types.NodeState) error {
	log.Printf("Pinging %s\n", n.Node)

	reply, err := app.PingNode(n.Client)
	if err != nil {
		return err
	}

	err = app.RefreshNode(n.Node)
	if err != nil {
		return err
	}

	for _, newNode := range reply.NodeList {
		lErr := app.AddNode(newNode.ID, newNode.Address)
		if lErr != nil {
			log.Println(err)
		}
	}

	return nil
}

// FirstPing is used on startup to ping nodes where only the address is known
func FirstPing(app *dnodes.DNodesApp, addr string) error {
	log.Printf("Pinging %s for the first time\n", addr)

	rpcClient, err := dnodes.CreateClient(addr)
	if err != nil {
		return err
	}

	reply, err := app.PingNode(rpcClient)
	if err != nil {
		return err
	}

	err = rpcClient.Close()
	if err != nil {
		return err
	}

	for _, newNode := range reply.NodeList {
		lErr := app.AddNode(newNode.ID, newNode.Address)
		if lErr != nil {
			log.Println(err)
		}
	}

	return nil
}
