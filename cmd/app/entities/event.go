package entities

import (
	"fmt"
	"time"
)

type Event struct {
	ID uint64 `json:"id" gorm:"primary_key;auto_increment"`

	Element   string `json:"element"`
	EventType string `json:"event_type" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Location  string `json:"location"`

	// user click positional data
	ClientX *int `json:"client_x"`
	ClientY *int `json:"client_y"`

	// evaluated on the backend
	SessionUUID string    `json:"-"`
	Username    string    `json:"-"`
	SourceIP    string    `json:"-"`
	Timestamp   time.Time `json:"-"`
}

type ClickData struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Value int `json:"value"`
}

func (e *Event) String() string {
	return fmt.Sprintf("session: %s, message: %s", e.SessionUUID, e.Message)
}
