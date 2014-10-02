package codec

import (
	"errors"
	"net"
	"reflect"

	"code.google.com/p/gogoprotobuf/proto"
)

var (
	ErrAlreadyRegistered = errors.New("Message already registered")
)

// CodecInterface describes the codec interface,
// which encodes/decodes protobuf messages from/to
// the TCP connection.
type CodecInterface interface {
	// Register registers a message so that the
	// codec can identify the message when reading
	// the TCP connection.
	Register(msg proto.Message) error
	// Encode encodes a message to bytes and
	// writes it to the TCP connection.
	Encode(msg proto.Message, conn *net.TCPConn) error
	// Decode reads bytes from the TCP connection
	// and decodes it to a message.
	Decode(msg proto.Message, conn *net.TCPConn) error
}

// ProtobufCodec implements the codec interface.
type ProtobufCodec struct {
	// registeredMessages is a map from message indices
	// to message types. The indices increase monotonically.
	registeredMessages map[int]reflect.Type
	// messageIndices is a map from message types
	// to message indices.
	messageIndices map[reflect.Type]int
}

// NewProtobufCodec creates and returns a ProtobufCodec.
func NewProtobufCodec() *ProtobufCodec {
	return &ProtobufCodec{
		registeredMessages: make(map[int]reflect.Type),
		messageIndices:     make(map[reflect.Type]int),
	}
}

// Register registers a message. Note this is not concurrent-safe.
func (pc *ProtobufCodec) Register(msg proto.Message) error {
	mtype := reflect.TypeOf(msg)
	if _, existed := pc.messageIndices[mtype]; existed {
		return ErrAlreadyRegistered
	}
	index := len(pc.messageIndices)
	pc.messageIndices[mtype] = index
	pc.registeredMessages[index] = mtype
	return nil
}
