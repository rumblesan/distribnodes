package dnodes

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/rumblesan/distribnodes/config"
	"github.com/rumblesan/distribnodes/types"
)

// RPCServer is the RPC handler
type RPCServer struct {
	App *DNodesApp
}

// Ping is sent to a node to check its liveness state and retrieve
// a list of all the other nodes it's aware of
type Ping struct {
	Node *types.DistribNode
}

// Pong is sent from a node to confirm its liveness state and return
// a list of all the other nodes it's aware of
type Pong struct {
	NodeList []*types.DistribNode
}

// PingNodeRPC handles a Ping RPC from a node
func (s *RPCServer) PingNodeRPC(pi Ping, po *Pong) error {
	log.Printf("PING from node %s\n", pi.Node.ID)
	*po = Pong{
		NodeList: s.App.GetFullNodeList(),
	}
	err := s.App.RefreshNode(pi.Node)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// PingNode sends a Ping RPC to a node
func (app *DNodesApp) PingNode(client *rpc.Client) (*Pong, error) {
	var reply Pong
	msg := Ping{
		Node: app.Self,
	}
	err := client.Call("RPCServer.PingNodeRPC", msg, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// StartRPCServer creates the server handler and starts listening
func StartRPCServer(app *DNodesApp, cfg *config.AppConfig) error {
	s := &RPCServer{App: app}
	err := rpc.Register(s)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", app.Self.Address)
	if err != nil {
		return err
	}

	log.Printf("Serving RPC server on %s", app.Self.Address)

	return http.Serve(listener, nil)
}

// CreateClient will create an rpc client for an address
func CreateClient(address string) (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", address)
}
