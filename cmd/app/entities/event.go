package entities

import (
	"fmt"
	"time"
)

type Event struct {
	Element   string `json:"element"`
	EventType string `json:"event_type" binding:"required"`
	Message   string `json:"message" binding:"required"`

	// evaluated on the backend
	SessionUUID string    `json:"-"`
	Username    string    `json:"-"`
	SourceIP    string    `json:"-"`
	Timestamp   time.Time `json:"-"`
}

func (e *Event) String() string {
	return fmt.Sprintf("session: %s, message: %s", e.SessionUUID, e.Message)
}
