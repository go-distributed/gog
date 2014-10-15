package codec

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/testify/assert"
)

const nsInOneSecond = 1000000000 // 10^9
var payload []byte

func genRandomMessage(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Int())
	}
	return b
}

func TestRegister(t *testing.T) {
	umsg := &message.UserMessage{
		Id:      proto.String("localhost:8080"),
		Payload: []byte("hello world"),
	}
	pc := NewProtobufCodec()
	pc.Register(umsg)
	assert.Panics(t, func() { pc.Register(umsg) })
}

func TestWriteMsgReadMsg(t *testing.T) {
	umsg1 := &message.UserMessage{
		Id:      proto.String("localhost:8080"),
		Payload: []byte("hello"),
	}
	umsg2 := &message.UserMessage{
		Id:      proto.String("localhost:8080"),
		Payload: []byte("world"),
	}
	pc := NewProtobufCodec()
	pc.Register(umsg1)
	rw := new(bytes.Buffer)
	assert.NoError(t, pc.WriteMsg(umsg1, rw))
	assert.NoError(t, pc.WriteMsg(umsg2, rw))
	msg1, err := pc.ReadMsg(rw)
	msg2, err := pc.ReadMsg(rw)
	assert.NoError(t, err)
	assert.Equal(t, umsg1, msg1)
	assert.Equal(t, umsg2, msg2)
}

func BenchmarkWriteMsgReadMsg(b *testing.B) {
	umsg := &message.UserMessage{
		Id:      proto.String("localhost:8080"),
		Payload: payload,
	}
	pc := NewProtobufCodec()
	pc.Register(umsg)
	rw := new(bytes.Buffer)

	for i := 0; i < b.N; i++ {
		assert.NoError(b, pc.WriteMsg(umsg, rw))
		_, err := pc.ReadMsg(rw)
		assert.NoError(b, err)
	}
}

func TestWriteMsgReadMsgThroughput(t *testing.T) {
	var result testing.BenchmarkResult

	fmt.Println("Small Message Payload (100 bytes)")
	payload = genRandomMessage(100)
	result = testing.Benchmark(BenchmarkWriteMsgReadMsg)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*10*10/1024/1024, "MB/s")

	fmt.Println("Medium Message Payload (1000 bytes)")
	payload = genRandomMessage(1000)
	result = testing.Benchmark(BenchmarkWriteMsgReadMsg)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*100*100/1024/1024, "MB/s")

	fmt.Println("Large Message Payload (1000*1000 bytes ~1m)")
	payload = genRandomMessage(1000 * 1000)
	result = testing.Benchmark(BenchmarkWriteMsgReadMsg)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*1000*1000/1024/1024, "MB/s")

}
