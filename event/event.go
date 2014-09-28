package event

import (
	"code.google.com/p/go-uuid/uuid"
)

// Event describes a non-message event (e.g. Connection lost).
type Event struct {
	// Etype is the event type.
	Etype int
	// Uuid indicates who triggers the event.
	Uuid  uuid.UUID
}
