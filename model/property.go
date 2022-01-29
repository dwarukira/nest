package model

import "github.com/google/uuid"

type Property struct {
	Base
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Owner       User      `json:"-"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"column:owner_id"`
	Units       *[]Unit   `json:"units"`
}

func (p Property) TableName() string {
	return "properties"
}
