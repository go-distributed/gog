package agent

import (
	"fmt"
	"net"
	"sync"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/go-distributed/gog/codec"
	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/message"
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
	// The id of the agent.
	id string
	// Configuration.
	cfg *config.Config
	// Active View.
	mua   sync.Mutex
	aView map[string]*node.Node
	// Passive View.
	mup   sync.Mutex
	pView map[string]*node.Node
	// TCP listener.
	ln *net.TCPListener
	// The codec.
	codec codec.CodecInterface
}

// NewAgent creates a new agent.
func (ag *Agent) NewAgent(cfg *config.Config) *Agent {
	// Create a codec and register messages.
	codec := codec.NewProtobufCodec()
	ag.codec.Register(&message.UserMessage{})
	ag.codec.Register(&message.Join{})
	ag.codec.Register(&message.ForwardJoin{})
	ag.codec.Register(&message.Disconnect{})
	ag.codec.Register(&message.Shuffle{})
	ag.codec.Register(&message.ShuffleReply{})

	return &Agent{
		cfg:   cfg,
		codec: codec,
		aView: make(map[string]*node.Node),
		pView: make(map[string]*node.Node),
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
	// TODO(yifan): Added a tick to trigger shuffle periodically.
	ag.serve()
	return nil
}

// serveNewConn listens on the TCP listener, waits for incoming connections.
func (ag *Agent) serve() {
	for {
		conn, err := ag.ln.AcceptTCP()
		if err != nil {
			log.Errorf("Agent.serve(): Failed to accept\n")
			continue
		}
		go ag.serveConn(conn)
	}
}

func (ag *Agent) serveConn(conn *net.TCPConn) {
	defer conn.Close()
	// TODO(Yifan): Set read time ount.

	for {
		msg, err := ag.codec.Decode(conn)
		if err != nil {
			log.Errorf("Agent.serveConn(): Failed to decode message: %v\n", err)
			// TODO(yifan): Now what? Drop the conn?
			// Update the view?
			return
		}
		// Dispatch messages.
		switch t := msg.(type) {
		case *message.Join:
			if !ag.handleJoin(conn, msg) {
				return
			}
		case *message.Neighbor:
			if !ag.handleNeighbor(conn, msg) {
				return
			}
		case *message.ForwardJoin:
			ag.handleForwardJoin(msg)
		case *message.Disconnect:
			ag.handleDisconnect(msg)
		case *message.Shuffle:
			ag.handleShuffle(msg)
		case *message.ShuffleReply:
			ag.handleShuffleReply(msg)
		case *message.UserMessage:
			ag.handleUserMessage(msg)
		default:
			log.Errorf("Agent.serveConn(): Unexpected message type: %T\n", t)
			// TODO(yifan): Now what? Drop the conn?
			// Update the view?
			return
		}
	}
}

// handleJoin() handles Join message, it returns true if it accepts and
// adds the node in the active view. As specified by the protocol. It should
// always accept Join requests, so it always returns true.
func (ag *Agent) handleJoin(conn *net.TCPConn, msg proto.Message) bool {
	fmt.Println("Fill me in")
	return true
}

// handleNeighbor() handles Neighbor message, it returns true if it accepts
// the request and adds the node in the active view. It returns false if it
// rejects the request.
func (ag *Agent) handleNeighbor(conn *net.TCPConn, msg proto.Message) bool {
	fmt.Println("Fill me in")
	return true
}

// handleForwardJoin() handles the ForwardJoin message, and decides whether
// it will add the original sender to the active view or passive view.
func (ag *Agent) handleForwardJoin(msg proto.Message) {
	fmt.Println("Fill me in")
	return
}

// handleDisconnect() handles Disconnect message. It will replace the node
// with another node from the passive view. And send Neighbor message to it.
func (ag *Agent) handleDisconnect(msg proto.Message) {
	fmt.Println("Fill me in")
	return
}

// handleShuffle() handles Shuffle message. It will send back a ShuffleReply
// message and update it's views.
func (ag *Agent) handleShuffle(msg proto.Message) {
	fmt.Println("Fill me in")
	return
}

// handleShuffleReply() handles ShuffleReply message. It will update it's views.
func (ag *Agent) handleShuffleReply(msg proto.Message) {
	fmt.Println("Fill me in")
	return
}

// handleUserMessage() handles user defined messages. It will forward the message
// to the nodes in its active view.
func (ag *Agent) handleUserMessage(msg proto.Message) {
	fmt.Println("Fill me in")
	return
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
