package agent

import (
	"fmt"
	"net"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/gogoprotobuf/proto"
	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/event"
	"github.com/go-distributed/gog/node"
	log "github.com/golang/glog"
)

// AgentInterface describes the interface of an agent.
type AgentInterface interface {
	// Serve starts a standalone agent, waiting for
	// incoming connections.
	Serve() error
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
	// The uuid of the agent.
	uuid uuid.UUID
	// Configuration.
	cfg *config.Config
	// Active View.
	aView []*node.Node
	// Passive View.
	pView []*node.Node
	// TCP listener.
	ln *net.TCPListener
	// Event channel.
	eventCh chan *event.Event
	// Message channel.
	messageCh chan proto.Message
}

// NewAgent creates a new agent.
func (ag *Agent) NewAgent(cfg *config.Config) *Agent {
	return &Agent{
		cfg:   cfg,
		aView: make([]*node.Node, 0, cfg.AViewSize),
		pView: make([]*node.Node, 0, cfg.PViewSize),
	}
}

// Serve starts a standalone agent, waiting for
// incoming connections.
func (ag *Agent) Serve() error {
	ln, err := net.ListenTCP(ag.cfg.Net, ag.cfg.LocalTCPAddr)
	if err != nil {
		return err
	}
	ag.ln = ln
	ag.serve()
	return nil
}

func (ag *Agent) serve() {
	go ag.serveNewConn()
	for {
		select {
		case evnt := <-ag.eventCh:
			ag.dispatchEvents(evnt)
		case msg := <-ag.messageCh:
			ag.dispatchMessages(msg)
		}
	}
}

// serveNewConn listens on the TCP listener, waits for incoming connections.
func (ag *Agent) serveNewConn() {
	for {
		conn, err := ag.ln.AcceptTCP()
		if err != nil {
			log.Errorf("Agent.serve(): Failed to accept")
			continue
		}
		go ag.handleNewConnection(conn)
	}
}

func (ag *Agent) dispatchEvents(evnt *event.Event) {
	switch evnt.Type {
	case event.NewConnectionEventType:
		ag.handleNewConnection(evnt.Content.(*net.TCPConn))
	}
}

func (ag *Agent) dispatchMessages(msg proto.Message) {

}

func (ag *Agent) handleNewConnection(conn *net.TCPConn) {
	
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
