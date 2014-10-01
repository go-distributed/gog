package event

import (
	"net"

	"code.google.com/p/go-uuid/uuid"
)

const (
	NewConnectionEventType = iota + 1
)

// Event describes a non-message event (e.g. Connection lost).
type Event struct {
	// Type is the event type.
	Type int
	// Uuid indicates who triggers the event.
	Uuid uuid.UUID
	// The event content.
	Content interface{}
}

func NewConnectionEvent(uuid uuid.UUID, conn *net.TCPConn) *Event {
	return &Event{NewConnectionEventType, uuid, conn}
}
