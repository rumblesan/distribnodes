package dnodes

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"

	"github.com/rumblesan/distribnodes/types"
	"github.com/segmentio/ksuid"
)

// NodeID is a type alias for string
type NodeID = string

// GenerateNodeID will generate a new ID for a node
func GenerateNodeID() NodeID {
	return "node_" + ksuid.New().String()
}

// DNodesApp defines the state of the application
type DNodesApp struct {
	Self     *types.DistribNode
	NodeList map[string]*types.NodeState
	AppMutex *sync.Mutex
}

// GetRemoteNodeList returns a list of all the remote nodes and their states
func (app *DNodesApp) GetRemoteNodeList() []*types.NodeState {
	app.AppMutex.Lock()
	defer app.AppMutex.Unlock()

	var nodes []*types.NodeState
	for _, ns := range app.NodeList {
		nodes = append(nodes, ns)
	}
	return nodes
}

// GetFullNodeList returns a node list including this node
func (app *DNodesApp) GetFullNodeList() []*types.DistribNode {
	app.AppMutex.Lock()
	defer app.AppMutex.Unlock()

	var nodes []*types.DistribNode
	for _, ns := range app.NodeList {
		nodes = append(nodes, ns.Node)
	}
	nodes = append(nodes, app.Self)
	return nodes
}

// AddNode adds a node to the NodeList
func (app *DNodesApp) AddNode(id string, address string) error {
	app.AppMutex.Lock()
	defer app.AppMutex.Unlock()

	_, prs := app.NodeList[id]
	if prs {
		log.Printf("Node %s already known about\n", id)
		return nil
	}
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return err
	}
	n := types.DistribNode{
		ID:      id,
		Address: address,
	}

	app.NodeList[n.ID] = &types.NodeState{
		Node:   &n,
		Seen:   time.Now(),
		Client: client,
	}
	return nil
}

// RefreshNode will update the seen at time for a node
func (app *DNodesApp) RefreshNode(n *types.DistribNode) error {
	app.AppMutex.Lock()
	defer app.AppMutex.Unlock()

	existing, prs := app.NodeList[n.ID]
	if !prs {
		rpcClient, err := rpc.DialHTTP("tcp", n.Address)
		if err != nil {
			return err
		}
		app.NodeList[n.ID] = &types.NodeState{
			Node:   n,
			Seen:   time.Now(),
			Client: rpcClient,
		}
	} else {
		existing.Seen = time.Now()
	}
	return nil
}

// RemoveNode removes a node from the NodeList
func (app *DNodesApp) RemoveNode(id string) error {
	app.AppMutex.Lock()
	defer app.AppMutex.Unlock()

	_, prs := app.NodeList[id]
	if !prs {
		return fmt.Errorf("no node with ID %s", id)
	}
	delete(app.NodeList, id)
	return nil
}
