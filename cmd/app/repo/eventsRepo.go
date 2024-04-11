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

func (repo EventsRepo) GetRecentEvents(limit int) ([]entities.Event, error) {
	var events []entities.Event

	err := repo.DB.Omit("ClientX", "ClientY", "SessionUUID").Limit(limit).Order("Timestamp DESC").Find(&events).Error

	return events, err
}
