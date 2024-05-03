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

func (repo EventsRepo) SaveMousePath(path entities.MousePath) error {
	return repo.Create(&path).Error
}

func (repo EventsRepo) GetRecentEvents(limit int) ([]entities.Event, error) {
	var events []entities.Event

	err := repo.DB.Omit("ClientX", "ClientY", "SessionUUID").Limit(limit).Order("Timestamp DESC").Find(&events).Error

	return events, err
}

func (repo EventsRepo) GetPathByFilter(session, location string) ([]entities.MousePath, error) {
	var paths []entities.MousePath

	err := repo.DB.
		Select([]string{"path", "session_uuid", "location"}).
		Where("location = ? AND session_uuid = ?", location, session).
		Find(&paths).
		Error

	return paths, err
}

func (repo EventsRepo) GetClicksByFilter(limit int, location string) ([]entities.ClickData, error) {
	var events []entities.Event
	var clicks = make([]entities.ClickData, 0, limit)

	err := repo.DB.
		Select("ClientX", "ClientY").
		Where("location = ? AND client_x IS NOT NULL AND client_y IS NOT NULL", location).
		Limit(limit).
		Find(&events).
		Error

	for _, e := range events {
		clicks = append(clicks, entities.ClickData{X: *e.ClientX, Y: *e.ClientY, Value: 1})
	}

	return clicks, err
}
