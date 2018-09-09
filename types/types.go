// Package types contains the shared types for the app
package types

import (
	"fmt"
	"net/rpc"
	"time"
)

// DistribNode is the shared information about given nodes in the network
type DistribNode struct {
	ID      string
	Address string
}

func (n DistribNode) String() string {
	return fmt.Sprintf("Node %s at %s", n.ID, n.Address)
}

// NodeState tracks the state of a remote DistribNode
type NodeState struct {
	Node   *DistribNode
	Seen   time.Time
	Client *rpc.Client
}
