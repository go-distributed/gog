package agent

import (
	"crypto/sha1"
	"fmt"
	"math/rand" // TODO(yifan): Need to change this??
	"net"
	"sync"
	"time"

	"github.com/go-distributed/gog/arraymap"
	"github.com/go-distributed/gog/codec"
	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/gog/node"

	"code.google.com/p/gogoprotobuf/proto"
	log "github.com/go-distributed/gog/log" // DEBUG
)

// Agent describes the interface of an agent.
type Agent interface {
	// Serve( starts a standalone agent, waiting for
	// incoming connections.
	Serve() error
	// Join joins the agent to the cluster.
	Join(addr string) error
	// Leave causes the agent to leave the cluster.
	Leave() error
	// Broadcast broadcasts a message to the cluster.
	Broadcast(msg []byte) error
	// Count does a broadcast and returns a channel of
	// nodes, which can be used to compute the broadcast
	// delay.
	Count(addr string) (chan *node.Node, error)
	// List prints the infomation in two views.
	List()
}

// agent implements the Agent interface.
type agent struct {
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
	codec codec.Codec
	// Message Buffer.
	msgBuffer *arraymap.ArrayMap
}

// NewAgent creates a new agent.
func NewAgent(cfg *config.Config) Agent {
	// Create a codec and register messages.
	codec := codec.NewProtobufCodec()
	codec.Register(&message.UserMessage{})
	codec.Register(&message.Join{})
	codec.Register(&message.ForwardJoin{})
	codec.Register(&message.Neighbor{})
	codec.Register(&message.NeighborReply{})
	codec.Register(&message.Disconnect{})
	codec.Register(&message.Shuffle{})
	codec.Register(&message.ShuffleReply{})

	return &agent{
		id:        cfg.AddrStr, // TODO(yifan): other id.
		cfg:       cfg,
		codec:     codec,
		aView:     make(map[string]*node.Node),
		pView:     make(map[string]*node.Node),
		msgBuffer: arraymap.NewArrayMap(),
	}
}

// Serve starts a standalone agent, waiting for
// incoming connections.
func (ag *agent) Serve() error {
	ln, err := net.ListenTCP(ag.cfg.Net, ag.cfg.LocalTCPAddr)
	if err != nil {
		log.Errorf("Serve() Cannot listen %v\n", err)
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
		// TODO(Yifan): Set read time ount.
		go ag.serveConn(conn)
	}
}

// serveConn() serves a connection.
func (ag *agent) serveConn(conn *net.TCPConn) {
	defer conn.Close()
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
			ag.handleJoin(conn, msg.(*message.Join))
		case *message.Neighbor:
			ag.handleNeighbor(conn, msg.(*message.Neighbor))
		case *message.ForwardJoin:
			ag.handleForwardJoin(msg.(*message.ForwardJoin))
		case *message.Disconnect:
			ag.handleDisconnect(msg.(*message.Disconnect))
		case *message.Shuffle:
			ag.handleShuffle(msg.(*message.Shuffle))
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

// chooseRandomNode() chooses a random node from the active view
// or passive view.
func chooseRandomNode(view map[string]*node.Node) *node.Node {
	index := rand.Intn(len(view))
	i := 0
	for _, node := range view {
		if i == index {
			return node
		}
		i++
	}
	panic("Must not get here!")
}

// addNodeActiveView() adds the node to the active view. If
// the active view is full, it will move one node from the active
// view to the passive view before adding the node.
// If the passive view is also full, it will drop a random node
// in the passive view.
func (ag *agent) addNodeActiveView(node *node.Node) {
	if node.Id == ag.id {
		return
	}
	if ag.aView[node.Id] != nil {
		return
	}
	if len(ag.aView) == ag.cfg.AViewSize {
		n := chooseRandomNode(ag.aView)
		ag.disconnect(n)
		delete(ag.aView, n.Id)
		ag.addNodePassiveView(n)
	}
	ag.aView[node.Id] = node
}

// addNodePassiveView() adds a node to the passive view. If
// the passive view is full, it will drop a random node.
func (ag *agent) addNodePassiveView(node *node.Node) {
	if node.Id == ag.id {
		return
	}
	if ag.aView[node.Id] != nil {
		return
	}
	if ag.pView[node.Id] != nil {
		return
	}
	if len(ag.pView) == ag.cfg.PViewSize {
		n := chooseRandomNode(ag.pView)
		delete(ag.pView, n.Id)
	}
	ag.pView[node.Id] = node
}

// handleJoin() handles Join message. If it accepts the request, it will add
// the node in the active view. As specified by the protocol, a node should
// always accept Join requests.
func (ag *agent) handleJoin(conn *net.TCPConn, msg *message.Join) {
	newNode := &node.Node{
		Id:   msg.GetId(),
		Addr: msg.GetAddr(),
		Conn: conn,
	}

	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()

	ag.addNodeActiveView(newNode)
	go ag.serveConn(newNode.Conn)

	// Send ForwardJoin message to all other the nodes in the active view.
	for _, node := range ag.aView {
		if node == newNode {
			continue
		}
		ag.forwardJoin(node, newNode, uint32(rand.Intn(ag.cfg.ARWL))) // TODO(yifan): go ag.forwardJoin()
	}
}

// handleNeighbor() handles Neighbor message. If the request is high priority,
// the receiver will always accept the request and add the node to its active view.
// If the request is low priority, then the request will only be accepted when
// there are empty slot in the active view.
func (ag *agent) handleNeighbor(conn *net.TCPConn, msg *message.Neighbor) {
	newNode := &node.Node{
		Id:   msg.GetId(),
		Addr: msg.GetAddr(),
		Conn: conn,
	}

	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()

	if len(ag.aView) == ag.cfg.AViewSize {
		if msg.GetPriority() == message.Neighbor_Low {
			ag.rejectNeighbor(newNode) // TODO(yifan): go ag.rejectNeighbor()
			// TODO(yifan): Add the node to passive view.
			return
		}
	}
	ag.addNodeActiveView(newNode)
	go ag.serveConn(newNode.Conn)
	ag.acceptNeighbor(newNode) // TODO(yifan): go ag.acceptNeighbor()
	return
}

// handleForwardJoin() handles the ForwardJoin message, and decides whether
// it will add the original sender to the active view or passive view.
func (ag *agent) handleForwardJoin(msg *message.ForwardJoin) {
	ttl := msg.GetTtl()
	newNode := &node.Node{
		Id:   msg.GetSourceId(),
		Addr: msg.GetSourceAddr(),
	}

	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()

	if ttl == 0 || len(ag.aView) == 1 { // TODO(yifan): Loose this?
		if ag.id != newNode.Id && ag.aView[newNode.Id] == nil {
			ag.neighbor(newNode, message.Neighbor_High)
		}
		return
	}
	if ttl == uint32(ag.cfg.PRWL) {
		ag.addNodePassiveView(newNode)
	}
	node := chooseRandomNode(ag.aView)
	ag.forwardJoin(node, newNode, ttl-1) // TODO(yifan): go ag.forwardJoin()
	return
}

// handleDisconnect() handles Disconnect message. It will replace the node
// with another node from the passive view. And send Neighbor message to it.
func (ag *agent) handleDisconnect(msg *message.Disconnect) {
	id := msg.GetId()

	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()

	node := ag.aView[id]
	if node == nil {
		return
	}
	delete(ag.aView, node.Id)
	node.Conn.Close()
	node.Conn = nil
	for len(ag.aView) == 0 { // We need at least one active node.
		n := chooseRandomNode(ag.pView)
		ag.neighbor(n, message.Neighbor_High)
	}
	ag.pView[id] = node
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
	ag.mua.Lock()
	defer ag.mua.Unlock()

	ms, err := time.ParseDuration("1ms")
	if err != nil {
		panic("failed to parse duration") // Shouldn't happen.
	}
	deadline := msg.GetTs() + ms.Nanoseconds()*int64(ag.cfg.MLife)
	hash := hashMessage(msg.GetPayload())
	if time.Now().UnixNano() >= deadline || ag.msgBuffer.Has(hash) {
		return
	}
	ag.msgBuffer.Append(hash, msg)
	for _, node := range ag.aView {
		ag.userMessage(node, msg) // TODO(yifan) go ag.userMessage
	}
	return
}

// Join joins the node to the cluster by contacting the nodes provied in the
// list.
func (ag *agent) Join(addr string) error {
	node := &node.Node{
		Id:   addr,
		Addr: addr,
	}
	tcpAddr, err := net.ResolveTCPAddr(ag.cfg.Net, node.Addr)
	if err != nil {
		// TODO(yifan) log.
		return err
	}
	conn, err := net.DialTCP(ag.cfg.Net, nil, tcpAddr)
	if err != nil {
		// TODO(yifan) log.
		return err
	}
	node.Conn = conn
	if err := ag.join(node); err != nil {
		return err
	}

	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()

	ag.addNodeActiveView(node)
	go ag.serveConn(node.Conn)
	return nil
}

// Leave causes the agent to leave the cluster.
func (ag *agent) Leave() error {
	return fmt.Errorf("Fill me in")
}

// Broadcast broadcasts a message to the cluster.
func (ag *agent) Broadcast(payload []byte) error {
	hash := hashMessage(payload)
	msg := &message.UserMessage{
		Id:      proto.String(ag.id),
		Payload: payload,
		Ts:      proto.Int64(time.Now().UnixNano()),
	}
	ag.msgBuffer.Append(hash, msg)
	ag.mua.Lock()
	defer ag.mua.Unlock()
	for _, node := range ag.aView {
		if err := ag.userMessage(node, msg); err != nil {
			// TODO: update view.
		}
	}
	return nil
}

// Count does a broadcast and returns a channel of
// nodes, which can be used to compute the broadcast
// delay.
func (ag *agent) Count(addr string) (chan *node.Node, error) {
	return nil, fmt.Errorf("Fill me in")
}

func (ag *agent) List() {
	ag.mua.Lock()
	ag.mup.Lock()
	defer ag.mua.Unlock()
	defer ag.mup.Unlock()
	fmt.Println("AView:")
	for _, node := range ag.aView {
		fmt.Println(node)
	}
	fmt.Println("PView:")
	for _, node := range ag.pView {
		fmt.Println(node)
	}

}

func hashMessage(msg []byte) [sha1.Size]byte {
	return sha1.Sum(msg)
}
