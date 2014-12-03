package node

import (
	"net"
)

// Node decribes a node in the overlay.
type Node struct {
	// Id is the node's identification.
	Id uint64 `json:"id"`
	// Addr is the network address of the node,
	// in the form of "host:port".
	Addr string `json:"address"`
	// Conn is the (TCP) connection to the node.
	// If the node is in the passive view, then the Conn could be
	// nil.
	Conn *net.TCPConn `json:"-"`
}
