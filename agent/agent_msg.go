package agent

import (
	"errors"
	"net"

	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/gog/node"

	"code.google.com/p/gogoprotobuf/proto"
)

var (
	ErrInvalidMessageType = errors.New("Invalid message type")
	ErrNoAvailablePeers   = errors.New("No available peers")
)

// disconnect() sends a Disconnect message to the node and close the connection.
// TODO(yifan): cache the connection.
func (ag *agent) disconnect(node *node.Node) {
	msg := &message.Disconnect{Id: proto.String(ag.id)}
	ag.codec.WriteMsg(msg, node.Conn) // TODO record err log.
	node.Conn.Close()
	node.Conn = nil
}

// forwardJoin() sends a ForwardJoin message to the node. The message
// will include the Id and Addr of the source node, as the receiver might
// use these information to establish a connection.
func (ag *agent) forwardJoin(node, newNode *node.Node, ttl uint32) {
	msg := &message.ForwardJoin{
		Id:         proto.String(ag.id),
		SourceId:   proto.String(newNode.Id),
		SourceAddr: proto.String(newNode.Addr),
		Ttl:        proto.Uint32(ttl),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}

// rejectNeighbor() sends a the NeighborReply with accept = false.
func (ag *agent) rejectNeighbor(node *node.Node) {
	msg := &message.NeighborReply{
		Id:     proto.String(ag.id),
		Accept: proto.Bool(false),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		// TODO log
	}
}

// acceptNeighbor() sends a the NeighborReply with accept = true.
func (ag *agent) acceptNeighbor(node *node.Node) {
	msg := &message.NeighborReply{
		Id:     proto.String(ag.id),
		Accept: proto.Bool(true),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}

func (ag *agent) join(node *node.Node) error {
	msg := &message.Join{
		Id:   proto.String(ag.id),
		Addr: proto.String(ag.cfg.AddrStr),
	}
	return ag.codec.WriteMsg(msg, node.Conn)
}

// neighbor() sends a Neighbor message, and wait for the reply.
// If the other side accepts the request, we add the node to the active view.
func (ag *agent) neighbor(node *node.Node, priority message.Neighbor_Priority) (error, bool) {
	if node.Conn == nil {
		addr, err := net.ResolveTCPAddr(ag.cfg.Net, node.Addr)
		if err != nil {
			// TODO(yifan) log.
			return err, false
		}
		conn, err := net.DialTCP(ag.cfg.Net, nil, addr)
		if err != nil {
			// TODO(yifan) log.
			return err, false
		}
		node.Conn = conn
	}
	msg := &message.Neighbor{
		Id:       proto.String(ag.id),
		Addr:     proto.String(ag.cfg.AddrStr),
		Priority: priority.Enum(),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		// TODO(yifan) log.
		node.Conn.Close()
		return err, false
	}
	recvMsg, err := ag.codec.ReadMsg(node.Conn)
	if err != nil {
		// TODO(yifan) log.
		node.Conn.Close()
		return err, false
	}
	reply, ok := recvMsg.(*message.NeighborReply)
	if !ok {
		node.Conn.Close()
		return ErrInvalidMessageType, false
	}
	if reply.GetAccept() {
		ag.addNodeActiveView(node)
		go ag.serveConn(node.Conn, node)
		return nil, true
	}
	ag.addNodePassiveView(node)
	return nil, false
}

// userMessage() sends a user message to the node.
func (ag *agent) userMessage(node *node.Node, msg proto.Message) {
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}

func (ag *agent) forwardShuffle(node *node.Node, msg *message.Shuffle) {
	msg.Id = proto.String(ag.id)
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}

func (ag *agent) shuffleReply(msg *message.Shuffle, candidates []*message.Candidate) error {
	// TODO use existing tcp.
	addr, err := net.ResolveTCPAddr(ag.cfg.Net, msg.GetAddr())
	if err != nil {
		// TODO(yifan) log
		return err
	}
	conn, err := net.DialTCP(ag.cfg.Net, nil, addr)
	if err != nil {
		return err
	}
	reply := &message.ShuffleReply{
		Id:         proto.String(ag.id),
		Candidates: candidates,
	}
	if err := ag.codec.WriteMsg(reply, conn); err != nil {
		// TODO log
		return err
	}
	conn.Close()
	return nil
}

func (ag *agent) shuffle(node *node.Node, candidates []*message.Candidate) {
	msg := &message.Shuffle{
		Id:         proto.String(ag.id),
		SourceId:   proto.String(ag.id),
		Addr:       proto.String(ag.cfg.AddrStr),
		Candidates: candidates,
		Ttl:        proto.Uint32(uint32(ag.cfg.SRWL)),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}
