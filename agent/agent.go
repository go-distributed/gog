package agent

import (
	"net"

	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/event"
	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/gog/node"
)

// Agent describes a gossip agent.
type Agent struct {
	// Configuration.
	cfg *config.Config

	// Active View.
	aView []*node.Node
	// Passive View.
	pView []*node.Node

	// TCP listener.
	ln *net.TCPListener
}
