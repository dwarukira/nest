package model

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type LeaseChargeTypeEmun string

const (
	SECURITY_DEPOSIT LeaseChargeTypeEmun = "SECURITY_DEPOSIT"
	OTHER_DEPOSIT    LeaseChargeTypeEmun = "OTHER_DEPOSIT"
	FEE              LeaseChargeTypeEmun = "FEE"
	OTHER            LeaseChargeTypeEmun = "OTHER"
	RENT             LeaseChargeTypeEmun = "RENT"
)

func (ct *LeaseChargeTypeEmun) Scan(value interface{}) error {
	*ct = LeaseChargeTypeEmun(value.(string))
	return nil
}

func (ct LeaseChargeTypeEmun) Value() (driver.Value, error) {
	return string(ct), nil
}

type LeaseChargeType struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}

type LeaseCharge struct {
	Base
	Name                 string               `gorm:"column:name" json:"name"`
	Description          string               `gorm:"column:description" json:"description"`
	Amount               int64                `gorm:"column:amount" json:"amount"`
	DueDate              time.Time            `gorm:"column:due_date" json:"due_date"`
	LeaseID              uuid.UUID            `gorm:"column:lease_id" json:"lease_id"`
	Lease                Lease                `json:"lease"`
	LeaseChargesPayments []LeaseChargePayment `json:"lease_charges_payment"`
	RemainingAmount      int64                `json:"remaining_amount" gorm:"->"`
	RecivedAmount        int64                `json:"recived_amount" gorm:"->"`
	LeaseChargeTypeID    uuid.UUID            `gorm:"column:lease_charge_type_id" json:"lease_charge_type_id"`
	ChargeType           LeaseChargeTypeEmun  `json:"charge_type" gorm:"charge_type" sql:"charge_type"`
}

func (u LeaseCharge) TableName() string {
	return "lease_charges"
}

type LeaseChargePayment struct {
	Base
	Amount        int64       `gorm:"column:amount" json:"amount"`
	PaymentDate   time.Time   `gorm:"column:payment_date" json:"payment_date"`
	LeaseChargeID uuid.UUID   `gorm:"column:lease_charge_id" json:"lease_charge_id"`
	LeaseCharge   LeaseCharge `json:"lease_charge"`
}

func (u LeaseChargePayment) TableName() string {
	return "lease_charges_payments"
}
