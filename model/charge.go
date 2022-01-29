package model

import (
	"time"

	"github.com/google/uuid"
)

type LeaseCharge struct {
	Base
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Amount      int64     `gorm:"column:amount" json:"amount"`
	DueDate     time.Time `gorm:"column:due_date" json:"due_date"`
	LeaseID     uuid.UUID `gorm:"column:lease_id" json:"lease_id"`
	Lease       Lease     `gorm:"column:lease" json:"lease"`
}

func (u LeaseCharge) TableName() string {
	return "lease_charges"
}
