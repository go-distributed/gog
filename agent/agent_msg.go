package agent

import (
	"math/rand"

	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/gog/node"

	"code.google.com/p/gogoprotobuf/proto"
)

// disconnect() sends a Disconnect message to the node, but doesn't close
// the connection.
func (ag *agent) disconnect(node *node.Node) error {
	return ag.codec.WriteMsg(
		&message.Disconnect{
			Id: proto.String(ag.id),
		},
		node.Conn)
}

// forwardJoin() sends a ForwardJoin message to the node. The message
// will include the Id and Addr of the source node, as the receiver might
// use these information to establish a connection.
func (ag *agent) forwardJoin(node, newNode *node.Node) error {
	return ag.codec.WriteMsg(
		&message.ForwardJoin{
			Id:         proto.String(ag.id),
			SourceId:   proto.String(node.Id),
			SourceAddr: proto.String(node.Addr),
			Ttl:        proto.Uint32(uint32(rand.Intn(ag.cfg.ARWL))),
		},
		node.Conn)
}

// rejectNeighbor() sends a the NeighborReply with accept = false.
func (ag *agent) rejectNeighbor(node *node.Node) error {
	return ag.codec.WriteMsg(
		&message.NeighborReply{
			Id:     proto.String(ag.id),
			Accept: proto.Bool(false),
		},
		node.Conn)
}

// acceptNeighbor() sends a the NeighborReply with accept = true.
func (ag *agent) acceptNeighbor(node *node.Node) error {
	return ag.codec.WriteMsg(
		&message.NeighborReply{
			Id:     proto.String(ag.id),
			Accept: proto.Bool(true),
		},
		node.Conn)
}
