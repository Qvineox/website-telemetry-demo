package repo

import (
	"gorm.io/gorm"
	"website-telemetry-demo/cmd/app/entities"
)

type UsersRepo struct {
	*gorm.DB
}

func NewUsersRepo(DB *gorm.DB) UsersRepo {
	return UsersRepo{DB: DB}
}

func (repo UsersRepo) GetAllUsernames() ([]string, error) {
	var events []entities.Event
	var usernames []string

	err := repo.DB.Select([]string{"username"}).Group("username").Find(&events).Error
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		usernames = append(usernames, event.Username)
	}

	return usernames, err
}

func (repo UsersRepo) GetUsernameSessions(username string) ([]string, error) {
	var events []entities.Event
	var sessions []string

	err := repo.DB.Select([]string{"session_uuid"}).Where("username = ?", username).Group("session_uuid").Find(&events).Error
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		sessions = append(sessions, event.SessionUUID)
	}

	return sessions, err
}
