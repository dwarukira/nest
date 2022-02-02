package model

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	Base
	FirstName      string     `json:"name" gorm:"column:first_name"`
	LastName       string     `json:"last_name" gorm:"column:last_name"`
	Email          string     `json:"email" gorm:"column:email"`
	PhoneNumber    string     `json:"phone" gorm:"column:phone_number"`
	InviteToken    string     `json:"invite_token" gorm:"column:invite_token"`
	InviteAccepted time.Time  `json:"invite_accepted" gorm:"column:invite_accepted"`
	InviteSent     time.Time  `json:"invite_sent" gorm:"column:invite_sent"`
	LeaseID        uuid.UUID  `json:"lease_id" gorm:"column:lease_id"`
	Lease          Lease      `json:"lease"`
	UserID         *uuid.UUID `json:"user_id" gorm:"column:user_id"`
	User           *User      `json:"user"`
	CreatedBy      *User      `json:"created_by" gorm:"column:created_by"`
	CreatedByID    *uuid.UUID `json:"created_by_id" gorm:"column:created_by_id"`
	Account        Account    `json:"account"`
	AccountID      *uuid.UUID `json:"account_id" gorm:"column:account_id"`
}

func (Tenant) TableName() string {
	return "tenants"
}
