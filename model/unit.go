package model

import "github.com/google/uuid"

type Unit struct {
	Base
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	DefaultRent int       `gorm:"column:default_rent" json:"default_rent"`
	Property    Property  `gorm:"foreignKey:id" json:"-"`
	PropertyID  uuid.UUID `gorm:"column:"property_id" json:"property_id"`
}

func (u Unit) TableName() string {
	return "units"
}
