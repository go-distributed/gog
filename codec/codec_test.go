package codec

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"code.google.com/p/go-uuid/uuid"
	"github.com/go-distributed/gog/message"
	"github.com/go-distributed/testify/assert"
)

const nsInOneSecond = 1000000000 // 10^9
var payload [][]byte

func genRandomMessage(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Int())
	}
	return b
}

func genRandomMessageGroup(num, length int) [][]byte {
	bb := make([][]byte, num)
	for i := range bb {
		bb[i] = genRandomMessage(length)
	}
	return bb
}

func TestRegister(t *testing.T) {
	umsg := &message.UserMessage{
		Uuid:    uuid.NewUUID(),
		Payload: [][]byte{[]byte("hello world")},
	}
	pc := NewProtobufCodec()
	assert.NoError(t, pc.Register(umsg))
	assert.Error(t, pc.Register(umsg))
}

func TestEncodeDecode(t *testing.T) {
	umsg1 := &message.UserMessage{
		Uuid:    uuid.NewUUID(),
		Payload: [][]byte{[]byte("hello")},
	}
	umsg2 := &message.UserMessage{
		Uuid:    uuid.NewUUID(),
		Payload: [][]byte{[]byte("world")},
	}
	pc := NewProtobufCodec()
	assert.NoError(t, pc.Register(umsg1))
	rw := new(bytes.Buffer)
	assert.NoError(t, pc.Encode(umsg1, rw))
	assert.NoError(t, pc.Encode(umsg2, rw))
	msg1, err := pc.Decode(rw)
	msg2, err := pc.Decode(rw)
	assert.NoError(t, err)
	assert.Equal(t, umsg1, msg1)
	assert.Equal(t, umsg2, msg2)
}

func BenchmarkEncodeDecode(b *testing.B) {
	umsg := &message.UserMessage{
		Uuid:    uuid.NewUUID(),
		Payload: payload,
	}
	pc := NewProtobufCodec()
	assert.NoError(b, pc.Register(umsg))
	rw := new(bytes.Buffer)

	for i := 0; i < b.N; i++ {
		assert.NoError(b, pc.Encode(umsg, rw))
		_, err := pc.Decode(rw)
		assert.NoError(b, err)
	}
}

func TestEncodeDecodeThroughput(t *testing.T) {
	var result testing.BenchmarkResult

	fmt.Println("Small Message (10 * 10 bytes)")
	payload = genRandomMessageGroup(10, 10)
	result = testing.Benchmark(BenchmarkEncodeDecode)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*10*10/1024/1024, "MB/s")

	fmt.Println("Medium Message (100 * 100 bytes)")
	payload = genRandomMessageGroup(100, 100)
	result = testing.Benchmark(BenchmarkEncodeDecode)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*100*100/1024/1024, "MB/s")

	fmt.Println("Large Message (1000 * 1000 bytes)")
	payload = genRandomMessageGroup(1000, 1000)
	result = testing.Benchmark(BenchmarkEncodeDecode)
	fmt.Println(result.N, "    ", result.NsPerOp(), "ns/op")
	fmt.Println("Throughput:", nsInOneSecond/result.NsPerOp()*1000*1000/1024/1024, "MB/s")

}
