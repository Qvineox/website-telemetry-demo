package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type MousePath struct {
	ID uint64 `json:"-" gorm:"primary_key;auto_increment;column:id"`

	Path     [][2]uint64 `json:"path" binding:"required" gorm:"-"`
	PathJSON string      `json:"-" gorm:"column:path"`

	Location string `json:"location" binding:"required"`

	SessionUUID string    `json:"-"`
	Username    string    `json:"-"`
	SourceIP    string    `json:"-"`
	Timestamp   time.Time `json:"-"`
}

func (p *MousePath) BeforeCreate(tx *gorm.DB) (err error) {
	pathBytes, err := json.Marshal(p.Path)
	if err != nil {
		return err
	}

	p.PathJSON = string(pathBytes)

	return
}

func (p *MousePath) AfterFind(tx *gorm.DB) (err error) {
	err = json.Unmarshal([]byte(p.PathJSON), &p.Path)
	if err != nil {
		return err
	}

	return
}
