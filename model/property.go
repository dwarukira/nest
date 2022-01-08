package model

import "github.com/google/uuid"

type Property struct {
	Base
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Owner       User      `gorm:"foreignKey:id" json:"-"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"column:owner_id"`
}

func (p Property) TableName() string {
	return "properties"
}
