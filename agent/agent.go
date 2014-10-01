package agent

import (
	"fmt"
	"net"

	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/node"
)

// AgentInterface describes the interface of an agent.
type AgentInterface interface {
	// Start starts a standalone agent, waiting for
	// incoming connections.
	Start() error
	// Join joins the agent to the cluster.
	Join(addr ...string) error
	// Leave causes the agent to leave the cluster.
	Leave() error
	// Broadcast broadcasts a message to the cluster.
	Broadcast(msg []byte) error
	// Count does a broadcast and returns a channel of
	// nodes, which can be used to compute the broadcast
	// delay.
	Count(addr string) (chan *node.Node, error)
}

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

// NewAgent creates a new agent.
func (ag *Agent) NewAgent(cfg *config.Config) *Agent {
	return &Agent{
		cfg:   cfg,
		aView: make([]*node.Node, cfg.AViewSize),
		pView: make([]*node.Node, cfg.PViewSize),
	}
}

// Start starts a standalone agent, waiting for
// incoming connections.
func (ag *Agent) Start() error {
	return fmt.Errorf("Fill me in")
}

// Leave causes the agent to leave the cluster.
func (ag *Agent) Leave() error {
	return fmt.Errorf("Fill me in")
}

// Broadcast broadcasts a message to the cluster.
func (ag *Agent) Broadcast(msg []byte) error {
	return fmt.Errorf("Fill me in")
}

// Count does a broadcast and returns a channel of
// nodes, which can be used to compute the broadcast
// delay.
func (ag *Agent) Count(addr string) (chan *node.Node, error) {
	return nil, fmt.Errorf("Fill me in")
}
