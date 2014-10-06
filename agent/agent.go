package agent

import (
	"fmt"
	"math/rand" // TODO(yifan): Need to change this??
	"net"
	"sync"

	"github.com/go-distributed/gog/codec"
	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/gog/node"

	"code.google.com/p/gogoprotobuf/proto"
	log "github.com/golang/glog"
)

// Agent describes the interface of an agent.
type Agent interface {
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

// agent implements the Agent interface.
type agent struct {
	// The id of the agent.
	id string
	// Configuration.
	cfg *config.Config
	// Active View.
	mua   sync.Mutex
	aView []*node.Node
	// Passive View.
	mup   sync.Mutex
	pView []*node.Node
	// TCP listener.
	ln *net.TCPListener
	// The codec.
	codec codec.Codec
}

// NewAgent creates a new agent.
func NewAgent(cfg *config.Config) Agent {
	// Create a codec and register messages.
	codec := codec.NewProtobufCodec()
	codec.Register(&message.UserMessage{})
	codec.Register(&message.Join{})
	codec.Register(&message.ForwardJoin{})
	codec.Register(&message.Disconnect{})
	codec.Register(&message.Shuffle{})
	codec.Register(&message.ShuffleReply{})

	return &agent{
		cfg:   cfg,
		codec: codec,
		aView: make([]*node.Node, 0, cfg.AViewSize),
		pView: make([]*node.Node, 0, cfg.PViewSize),
	}
}

// Serve starts a standalone agent, waiting for
// incoming connections.
func (ag *agent) Serve() error {
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
func (ag *agent) serve() {
	for {
		conn, err := ag.ln.AcceptTCP()
		if err != nil {
			log.Errorf("Agent.serve(): Failed to accept\n")
			continue
		}
		go ag.serveConn(conn)
	}
}

func (ag *agent) serveConn(conn *net.TCPConn) {
	defer conn.Close()
	// TODO(Yifan): Set read time ount.

	for {
		msg, err := ag.codec.ReadMsg(conn)
		if err != nil {
			log.Errorf("Agent.serveConn(): Failed to decode message: %v\n", err)
			// TODO(yifan): Now what? Drop the conn?
			// Update the view?
			return
		}
		// Dispatch messages.
		switch t := msg.(type) {
		case *message.Join:
			if !ag.handleJoin(conn, msg.(*message.Join)) {
				return
			}
		case *message.Neighbor:
			if !ag.handleNeighbor(conn, msg.(*message.Neighbor)) {
				return
			}
		case *message.ForwardJoin:
			ag.handleForwardJoin(msg.(*message.ForwardJoin))
		case *message.Disconnect:
			ag.handleDisconnect(msg.(*message.Disconnect))
		case *message.Shuffle:
			ag.handleShuffle(msg.(*message.Shuffle))
		case *message.ShuffleReply:
			ag.handleShuffleReply(msg.(*message.ShuffleReply))
		case *message.UserMessage:
			ag.handleUserMessage(msg.(*message.UserMessage))
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
func (ag *agent) handleJoin(conn *net.TCPConn, msg *message.Join) bool {
	srcId, srcAddr := msg.GetId(), conn.RemoteAddr().String()

	ag.mua.Lock()
	defer ag.mua.Unlock()

	index := len(ag.aView)
	if index == ag.cfg.AViewSize {
		// Choose a victim, send Disconnect.
		index = rand.Intn(ag.cfg.AViewSize)
		ag.disconnect(ag.aView[index])
	}
	// Add the node to our active view.
	ag.aView[index] = &node.Node{
		Id:   srcId,
		Addr: srcAddr,
		Conn: conn,
	}
	// Send ForwardJoin message to all other the nodes in the active view.
	for i, node := range ag.aView {
		if i == index {
			continue
		}
		if err := ag.forwardJoin(node, srcId, srcAddr); err != nil {
			// TODO(yifan): Check error and replace the node.
		}
	}
	return true
}

// disconnect() sends a Disconnect message to the node, it does not
// return an error because when we send Disconnet, we are going to
// close the connection anyway.
func (ag *agent) disconnect(node *node.Node) {
	ag.codec.WriteMsg(
		&message.Disconnect{
			Id: proto.String(ag.id),
		},
		node.Conn)
}

// forwardJoin() sends a ForwardJoin message to the node. The message
// will include the srcId and srcAddr provided, as the receiver might
// use these information to establish a connection.
func (ag *agent) forwardJoin(node *node.Node, srcId, srcAddr string) error {
	return ag.codec.WriteMsg(
		&message.ForwardJoin{
			Id:         proto.String(ag.id),
			SourceId:   proto.String(srcId),
			SourceAddr: proto.String(srcAddr),
			Ttl:        proto.Uint32(uint32(rand.Intn(ag.cfg.ARWL))),
		},
		node.Conn)
}

// handleNeighbor() handles Neighbor message, it returns true if it accepts
// the request and adds the node in the active view. It returns false if it
// rejects the request.
func (ag *agent) handleNeighbor(conn *net.TCPConn, msg *message.Neighbor) bool {
	fmt.Println("Fill me in")
	return true
}

// handleForwardJoin() handles the ForwardJoin message, and decides whether
// it will add the original sender to the active view or passive view.
func (ag *agent) handleForwardJoin(msg *message.ForwardJoin) {
	fmt.Println("Fill me in")
	return
}

// handleDisconnect() handles Disconnect message. It will replace the node
// with another node from the passive view. And send Neighbor message to it.
func (ag *agent) handleDisconnect(msg *message.Disconnect) {
	fmt.Println("Fill me in")
	return
}

// handleShuffle() handles Shuffle message. It will send back a ShuffleReply
// message and update it's views.
func (ag *agent) handleShuffle(msg *message.Shuffle) {
	fmt.Println("Fill me in")
	return
}

// handleShuffleReply() handles ShuffleReply message. It will update it's views.
func (ag *agent) handleShuffleReply(msg *message.ShuffleReply) {
	fmt.Println("Fill me in")
	return
}

// handleUserMessage() handles user defined messages. It will forward the message
// to the nodes in its active view.
func (ag *agent) handleUserMessage(msg *message.UserMessage) {
	fmt.Println("Fill me in")
	return
}

// Join joins the node to the cluster by contacting the nodes provied in the
// list.
func (ag *agent) Join(addr ...string) error {
	return fmt.Errorf("Fill me in")
}

// Leave causes the agent to leave the cluster.
func (ag *agent) Leave() error {
	return fmt.Errorf("Fill me in")
}

// Broadcast broadcasts a message to the cluster.
func (ag *agent) Broadcast(msg []byte) error {
	return fmt.Errorf("Fill me in")
}

// Count does a broadcast and returns a channel of
// nodes, which can be used to compute the broadcast
// delay.
func (ag *agent) Count(addr string) (chan *node.Node, error) {
	return nil, fmt.Errorf("Fill me in")
}
