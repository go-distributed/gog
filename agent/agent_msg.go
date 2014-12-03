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
	msg := &message.Disconnect{Id: proto.Uint64(ag.id)}
	ag.codec.WriteMsg(msg, node.Conn) // TODO record err log.
	node.Conn.Close()
}

// forwardJoin() sends a ForwardJoin message to the node. The message
// will include the Id and Addr of the source node, as the receiver might
// use these information to establish a connection.
func (ag *agent) forwardJoin(node, newNode *node.Node, ttl uint32) {
	msg := &message.ForwardJoin{
		Id:         proto.Uint64(ag.id),
		SourceId:   proto.Uint64(newNode.Id),
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
		Id:     proto.Uint64(ag.id),
		Accept: proto.Bool(false),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		// TODO log
	}
}

// acceptNeighbor() sends a the NeighborReply with accept = true.
func (ag *agent) acceptNeighbor(node *node.Node) {
	msg := &message.NeighborReply{
		Id:     proto.Uint64(ag.id),
		Accept: proto.Bool(true),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}

func (ag *agent) join(node *node.Node) (bool, error) {
	msg := &message.Join{
		Id:   proto.Uint64(ag.id),
		Addr: proto.String(ag.cfg.AddrStr),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
		return false, err
	}
	recvMsg, err := ag.codec.ReadMsg(node.Conn)
	if err != nil {
		// TODO(yifan) log.
		node.Conn.Close()
		return false, err
	}
	reply, ok := recvMsg.(*message.JoinReply)
	if !ok {
		node.Conn.Close()
		return false, ErrInvalidMessageType
	}
	id, accepted := reply.GetId(), reply.GetAccept()
	node.Id = id

	if accepted {
		ag.aView.Lock()
		ag.pView.Lock()
		defer ag.aView.Unlock()
		defer ag.pView.Unlock()

		ag.addNodeActiveView(node)
		go ag.serveConn(node.Conn, node)
	}
	return accepted, nil
}

func (ag *agent) acceptJoin(node *node.Node) error {
	msg := &message.JoinReply{
		Id:     proto.Uint64(ag.id),
		Accept: proto.Bool(true),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
		return err
	}
	return nil
}

func (ag *agent) rejectJoin(node *node.Node) error {
	defer node.Conn.Close()
	msg := &message.JoinReply{
		Id:     proto.Uint64(ag.id),
		Accept: proto.Bool(false),
	}

	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		return err
	}
	return nil
}

// neighbor() sends a Neighbor message, and wait for the reply.
// If the other side accepts the request, we add the node to the active view.
func (ag *agent) neighbor(node *node.Node, priority message.Neighbor_Priority) (bool, error) {
	ag.aView.Unlock()
	ag.pView.Unlock()
	defer ag.aView.Lock()
	defer ag.pView.Lock()

	addr, err := net.ResolveTCPAddr(ag.cfg.Net, node.Addr)
	if err != nil {
		// TODO(yifan) log.
		return false, err
	}
	conn, err := net.DialTCP(ag.cfg.Net, nil, addr)
	if err != nil {
		// TODO(yifan) log.
		return false, err
	}
	node.Conn = conn

	msg := &message.Neighbor{
		Id:       proto.Uint64(ag.id),
		Addr:     proto.String(ag.cfg.AddrStr),
		Priority: priority.Enum(),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		// TODO(yifan) log.
		node.Conn.Close()
		return false, err
	}
	recvMsg, err := ag.codec.ReadMsg(node.Conn)
	if err != nil {
		// TODO(yifan) log.
		node.Conn.Close()
		return false, err
	}
	reply, ok := recvMsg.(*message.NeighborReply)
	if !ok {
		node.Conn.Close()
		return false, ErrInvalidMessageType
	}

	accepted := reply.GetAccept()

	ag.aView.Lock()
	ag.pView.Lock()
	if accepted {
		ag.addNodeActiveView(node)
		go ag.serveConn(node.Conn, node)
	} else {
		ag.addNodePassiveView(node)
	}
	ag.aView.Unlock()
	ag.pView.Unlock()

	return accepted, nil
}

// userMessage() sends a user message to the node.
func (ag *agent) userMessage(node *node.Node, msg proto.Message) {
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		// Record this message, so we can resend it later.
		umsg := msg.(*message.UserMessage)
		hash := hashMessage(umsg.GetPayload())

		ag.failmsgBuffer.Lock()
		ag.failmsgBuffer.Append(hash, msg)
		ag.failmsgBuffer.Unlock()

		node.Conn.Close()
	}
}

func (ag *agent) forwardShuffle(node *node.Node, msg *message.Shuffle) {
	msg.Id = proto.Uint64(ag.id)
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
		Id:         proto.Uint64(ag.id),
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
		Id:         proto.Uint64(ag.id),
		SourceId:   proto.Uint64(ag.id),
		Addr:       proto.String(ag.cfg.AddrStr),
		Candidates: candidates,
		Ttl:        proto.Uint32(uint32(ag.cfg.SRWL)),
	}
	if err := ag.codec.WriteMsg(msg, node.Conn); err != nil {
		node.Conn.Close()
	}
}
