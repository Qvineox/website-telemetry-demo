package entities

import "time"

type Comment struct {
	Author string `json:"author"`
	Text   string `json:"text"`

	LessonID *uint64 `json:"lesson_id"`

	CreatedAt time.Time `json:"created_at"`
}
