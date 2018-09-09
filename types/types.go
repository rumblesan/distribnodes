// Package types contains the shared types for the app
package types

import (
	"fmt"
	"time"
)

// DistribNode is the shared information about given nodes in the network
type DistribNode struct {
	ID      string
	Address string
}

// Copy makes a copy of a DistribNode
func (n DistribNode) Copy() DistribNode {
	return DistribNode{ID: n.ID, Address: n.Address}
}

func (n DistribNode) String() string {
	return fmt.Sprintf("Node %s at %s", n.ID, n.Address)
}

// NodeState tracks the state of a remote DistribNode
type NodeState struct {
	Node DistribNode
	Seen time.Time
}

// AppConfig holds the configuration for the application
type AppConfig struct {
	PingDuration time.Duration
}
