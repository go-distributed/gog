package node

import (
	"code.google.com/p/go-uuid/uuid"
	"net"
)

// Node decribes a node in the overlay.
type Node struct {
	// Uuid is the node's identification.
	Uuid uuid.UUID
	// Addr is the network address of the node,
	// it the form of "host:port".
	Addr string
	// Conn is the (TCP) connection to the node.
	// If the node is in the passive view, then the Conn could be
	// nil.
	Conn *net.TCPConn
}
