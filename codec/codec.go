package codec

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"reflect"

	"code.google.com/p/gogoprotobuf/proto"
)

const sizeOfUint8 = 1

var (
	ErrMessageAlreadyRegistered = errors.New("Message already registered")
	ErrMessageNotRegistered     = errors.New("Message not registered")
	ErrCannotWriteMessage       = errors.New("Cannot write message")
)

// CodecInterface describes the codec interface,
// which encodes/decodes protobuf messages from/to
// an io.Reader/Writer
type CodecInterface interface {
	// Register registers a message so that the
	// codec can identify the message when reading
	// the TCP connection.
	Register(msg proto.Message) error
	// Encode encodes a message to bytes and
	// writes it to the io.Writer
	Encode(msg proto.Message, w io.Writer) error
	// Decode reads bytes from the io.Reader
	// and decodes it to a message.
	Decode(r io.Reader) (proto.Message, error)
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
func (pc *ProtobufCodec) Register(msg proto.Message) error {
	mtype := reflect.TypeOf(msg)
	if _, existed := pc.messageIndices[mtype]; existed {
		return ErrMessageAlreadyRegistered
	}
	index := uint8(len(pc.messageIndices))
	pc.messageIndices[mtype] = index
	pc.registeredMessages[index] = mtype
	return nil
}

// Encode encodes a message to bytes and writes it to the io.Writer.
func (pc *ProtobufCodec) Encode(msg proto.Message, w io.Writer) error {
	index, existed := pc.messageIndices[reflect.TypeOf(msg)]
	if !existed {
		return ErrMessageNotRegistered
	}
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(w)

	// 1. Write the length.
	if err := binary.Write(bw, binary.LittleEndian, int32(len(b))); err != nil {
		return err
	}
	// 2. Write the type.
	if err := binary.Write(bw, binary.LittleEndian, index); err != nil {
		return err
	}
	// 3. Write the bytes.
	n, err := bw.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return ErrCannotWriteMessage
	}
	return bw.Flush()
}

// Decode reads bytes from an io.Reader and decode it to a message.
func (pc *ProtobufCodec) Decode(r io.Reader) (proto.Message, error) {
	var length int32
	var index uint8

	// 1. Read the length.
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	br := bufio.NewReaderSize(r, int(length)+sizeOfUint8)
	// 1. Read the type.
	if err := binary.Read(br, binary.LittleEndian, &index); err != nil {
		return nil, err
	}
	// 3. Read the bytes.
	b := make([]byte, length)
	_, err := io.ReadFull(br, b)
	if err != nil {
		return nil, err
	}
	// Decode.
	mtype, existed := pc.registeredMessages[index]
	if !existed {
		return nil, ErrMessageNotRegistered
	}
	msg := reflect.New(mtype.Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(b, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
