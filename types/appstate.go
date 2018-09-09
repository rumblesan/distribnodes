package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

// NodeID is a type alias for string
type NodeID = string

// GenerateNodeID will generate a new ID for a node
func GenerateNodeID() NodeID {
	return "node_" + ksuid.New().String()
}

// AppState defines the state of the application
type AppState struct {
	Self     DistribNode
	NodeList map[string]NodeState
	AppMutex *sync.Mutex
}

// GetNodeList returns the node list from the application state
func (as *AppState) GetNodeList() []DistribNode {
	as.AppMutex.Lock()
	defer as.AppMutex.Unlock()

	var nodes []DistribNode
	for _, ns := range as.NodeList {
		nodes = append(nodes, ns.Node.Copy())
	}
	return nodes
}

// AddNode adds a node to the NodeList
func (as *AppState) AddNode(n DistribNode) {
	as.AppMutex.Lock()
	defer as.AppMutex.Unlock()

	as.NodeList[n.ID] = NodeState{
		Node: n,
		Seen: time.Now(),
	}
}

// RemoveNode removes a node from the NodeList
func (as *AppState) RemoveNode(id string) error {
	as.AppMutex.Lock()
	defer as.AppMutex.Unlock()

	_, prs := as.NodeList[id]
	if !prs {
		return fmt.Errorf("no node with ID %s", id)
	}
	delete(as.NodeList, id)
	return nil
}

// HandshakeInfo creates the data to send in a handshake request
func (as *AppState) HandshakeInfo() []DistribNode {
	as.AppMutex.Lock()
	defer as.AppMutex.Unlock()

	var nodes []DistribNode
	for _, ns := range as.NodeList {
		nodes = append(nodes, ns.Node.Copy())
	}
	nodes = append(nodes, as.Self.Copy())
	return nodes
}
