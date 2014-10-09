package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"reflect"

	"code.google.com/p/gogoprotobuf/proto"

	log "github.com/go-distributed/gog/log" // DEBUG
)

const (
	sizeOfUint8 = 1
	sizeOfInt32 = 4
)

var (
	ErrMessageAlreadyRegistered = errors.New("Message already registered")
	ErrMessageNotRegistered     = errors.New("Message not registered")
	ErrCannotWriteMessage       = errors.New("Cannot write message")
)

// Codec describes the codec interface,
// which encodes/decodes protobuf messages from/to
// an io.Reader/Writer
type Codec interface {
	// Register registers a message so that the
	// codec can identify the message when reading
	// the TCP connection.
	Register(msg proto.Message)
	// WriteMsg encodes a message to bytes and
	// writes it to the io.Writer
	WriteMsg(msg proto.Message, w io.Writer) error
	// ReadMsg reads bytes from the io.Reader
	// and decodes it to a message.
	ReadMsg(r io.Reader) (proto.Message, error)
}

// ProtobufCodec implements the codec interface.
type ProtobufCodec struct {
	// registeredMessages is a map from message indices
	// to message types. The indices increase monotonically.
	registeredMessages map[uint8]reflect.Type
	// messageIndices is a map from message types
	// to message indices.
	messageIndices map[reflect.Type]uint8
}

// NewProtobufCodec creates and returns a ProtobufCodec.
func NewProtobufCodec() *ProtobufCodec {
	return &ProtobufCodec{
		registeredMessages: make(map[uint8]reflect.Type),
		messageIndices:     make(map[reflect.Type]uint8),
	}
}

// Register registers a message. Note this is not concurrent-safe.
func (pc *ProtobufCodec) Register(msg proto.Message) {
	mtype := reflect.TypeOf(msg)
	if _, existed := pc.messageIndices[mtype]; existed {
		panic("Message already registered")
	}
	index := uint8(len(pc.messageIndices))
	pc.messageIndices[mtype] = index
	pc.registeredMessages[index] = mtype
	return
}

// WriteMsg encodes a message to bytes and writes it to the io.Writer.
func (pc *ProtobufCodec) WriteMsg(msg proto.Message, w io.Writer) error {
	log.Infoln("Send:", msg, "to:", w.(*net.TCPConn).RemoteAddr())
	index, existed := pc.messageIndices[reflect.TypeOf(msg)]
	if !existed {
		return ErrMessageNotRegistered
	}
	buf := new(bytes.Buffer)

	// Encode.
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	// Write the length.
	if err := binary.Write(buf, binary.LittleEndian, int32(len(b)+sizeOfUint8)); err != nil {
		return err
	}
	// Write the type.
	if err := binary.Write(buf, binary.LittleEndian, index); err != nil {
		return err
	}
	// Write the bytes.
	buf.Write(b)
	if _, err = buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

// ReadMsg reads bytes from an io.Reader and decode it to a message.
func (pc *ProtobufCodec) ReadMsg(r io.Reader) (proto.Message, error) {
	var length int32

	// Read the length.
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	b := make([]byte, length)
	// Read the type and bytes.
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	// Get the index.
	index := uint8(b[0])
	// Decode.
	mtype, existed := pc.registeredMessages[index]
	if !existed {
		return nil, ErrMessageNotRegistered
	}
	msg := reflect.New(mtype.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(b[1:], msg); err != nil {
		return nil, err
	}
	log.Infoln("Recv:", msg, "from:", r.(*net.TCPConn).RemoteAddr())
	return msg, nil
}
