package repo

import (
	"gorm.io/gorm"
	"website-telemetry-demo/cmd/app/entities"
)

type EventsRepo struct {
	*gorm.DB
}

func NewEventsRepo(DB *gorm.DB) EventsRepo {
	return EventsRepo{DB: DB}
}

func (repo EventsRepo) SaveEvent(event entities.Event) error {
	return repo.Create(&event).Error
}
