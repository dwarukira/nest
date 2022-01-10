package model

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type LeaseStatusType string

const (
	DRAFT    LeaseStatusType = "DRAFT"
	ACTIVE   LeaseStatusType = "ACTIVE"
	INACTIVE LeaseStatusType = "INACTIVE"
)

func (ct *LeaseStatusType) Scan(value interface{}) error {
	*ct = LeaseStatusType(value.(string))
	return nil
}

func (ct LeaseStatusType) Value() (driver.Value, error) {
	return string(ct), nil
}

type Lease struct {
	Base
	LeaseNumber       string          `gorm:"column:lease_number" json:"lease_number"`
	StartDate         time.Time       `gorm:"column:start_date" json:"start_date"`
	EndDate           time.Time       `gorm:"column:end_date" json:"end_date"`
	MonthlyRent       int             `gorm:"column:monthly_rent" json:"rent"`
	SecurityDeposit   int             `gorm:"column:security_deposit" json:"security_deposit"`
	UnitID            uuid.UUID       `gorm:"column:unit_id" json:"unit_id"`
	LeaseStatus       LeaseStatusType `sql:"status" gorm:"column:status" json:"lease_status"`
	RentDueDayOfMonth int             `gorm:"column:rent_due_day_of_month" json:"rent_due_day_of_month"`
	Unit              Unit            `json:"unit"`
	Tenants           *[]Tenant       `json:"tenants"`
}

func (l *Lease) TableName() string {
	return "leases"
}
