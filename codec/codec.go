package codec

import (
	"net"

	"code.google.com/p/gogoprotobuf/proto"
)

type msgType int

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
