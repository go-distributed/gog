package arraymap

import (
	"testing"

	"github.com/go-distributed/testify/assert"
)

func TestSimple(t *testing.T) {
	am := NewArrayMap()
	for i := 0; i < 42; i++ {
		am.Append("foo", "bar")
	}
	assert.Equal(t, 1, am.Len())
	am.Append("bar", "foo")
	assert.Equal(t, 2, am.Len())

	am.Append("hello", "world")
	assert.Equal(t, 3, am.Len())

	am.Append("world", "hello")
	assert.Equal(t, 4, am.Len())

	assert.Equal(t, "foo", am.GetKeyAt(0))
	assert.Equal(t, "bar", am.GetKeyAt(1))
	assert.Equal(t, "hello", am.GetKeyAt(2))
	assert.Equal(t, "world", am.GetKeyAt(3))

	assert.Equal(t, "bar", am.GetValueAt(0))
	assert.Equal(t, "foo", am.GetValueAt(1))
	assert.Equal(t, "world", am.GetValueAt(2))
	assert.Equal(t, "hello", am.GetValueAt(3))

	assert.True(t, am.Has("foo"))
	assert.True(t, am.Has("bar"))
	assert.True(t, am.Has("hello"))
	assert.True(t, am.Has("world"))

	assert.False(t, am.Has("good"))

	assert.Panics(t, func() { am.RemoveAt(10) })
	am.RemoveAt(2)
	assert.False(t, am.Has("hello"))
	assert.Equal(t, 3, am.Len())
	assert.Equal(t, 3, len(am.positions))
	assert.Equal(t, 3, len(am.keys))
	assert.Equal(t, 3, len(am.values))

	assert.Equal(t, "world", am.GetKeyAt(2))
	assert.Equal(t, "hello", am.GetValueAt(2))
	am.RemoveAt(0)
	assert.Equal(t, "world", am.GetKeyAt(0))
	assert.Equal(t, "hello", am.GetValueAt(0))
	am.RemoveAt(0)
	assert.Equal(t, "bar", am.GetKeyAt(0))
	assert.Equal(t, "foo", am.GetValueAt(0))
	am.Remove("bar")
	assert.Panics(t, func() { am.RemoveAt(0) })
	am.Remove("aaaa")

	assert.Equal(t, 0, am.Len())
	assert.Equal(t, 0, len(am.positions))
	assert.Equal(t, 0, len(am.keys))
	assert.Equal(t, 0, len(am.values))
}
