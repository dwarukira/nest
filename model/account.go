package model

import (
	"database/sql/driver"
	"errors"

	"github.com/go-playground/validator/v10"
)

// AccountStatus represents the status of an account.
type AccountStatus string

// AccountStatus values define the status field of a user account.
const (
	// AccountStatus_Active defines the state when a user can access an account.
	AccountStatus_Active AccountStatus = "active"
	// AccountStatus_Pending defined the state when an account was created but
	// not activated.
	AccountStatus_Pending AccountStatus = "pending"
	// AccountStatus_Disabled defines the state when a user has been disabled from
	// accessing an account.
	AccountStatus_Disabled AccountStatus = "disabled"
)

// AccountStatus_Values provides list of valid AccountStatus values.
var AccountStatus_Values = []AccountStatus{
	AccountStatus_Active,
	AccountStatus_Pending,
	AccountStatus_Disabled,
}

// type Account struct {
// 	Base
// 	Name          string          `json:"name" validate:"required,unique" example:"Company Name"`
// 	Address1      string          `json:"address1" validate:"required" example:"221 Tatitlek Ave"`
// 	Address2      string          `json:"address2" validate:"omitempty" example:"Box #1832"`
// 	City          string          `json:"city" validate:"required" example:"Valdez"`
// 	Region        string          `json:"region" validate:"required" example:"AK"`
// 	Country       string          `json:"country" validate:"required" example:"USA"`
// 	Zipcode       string          `json:"zipcode" validate:"required" example:"99686"`
// 	Status        AccountStatus   `json:"status" validate:"omitempty,oneof=active pending disabled" swaggertype:"string" enums:"active,pending,disabled" example:"active"`
// 	Timezone      string          `json:"timezone" validate:"omitempty" example:"America/Anchorage"`
// 	SignupUserID  *sql.NullString `json:"signup_user_id,omitempty" validate:"omitempty,uuid" swaggertype:"string" example:"d69bdef7-173f-4d29-b52c-3edc60baf6a2"`
// 	BillingUserID *sql.NullString `json:"billing_user_id,omitempty" validate:"omitempty,uuid" swaggertype:"string" example:"d69bdef7-173f-4d29-b52c-3edc60baf6a2"`
// 	CreatedAt     time.Time       `json:"created_at"`
// 	UpdatedAt     time.Time       `json:"updated_at"`
// 	ArchivedAt    *pq.NullTime    `json:"archived_at,omitempty"`
// }

// Scan supports reading the AccountStatus value from the database.
func (s *AccountStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*s = AccountStatus(string(asBytes))
	return nil
}

// Value converts the AccountStatus value to be stored in the database.
func (s AccountStatus) Value() (driver.Value, error) {
	v := validator.New()

	errs := v.Var(s, "required,oneof=active invited disabled")
	if errs != nil {
		return nil, errs
	}

	return string(s), nil
}

// String converts the AccountStatus value to a string.
func (s AccountStatus) String() string {
	return string(s)
}

type Account struct {
	Base
	Name    string    `json:"name" example:"Company Name"`
	Users   []*User   `json:"users" gorm:"many2many:memberships;"`
	Tenants []*Tenant `json:"tenants"`
}

func (Account) TableName() string {
	return "accounts"
}

type Membership struct {
	Base
	AccountID uint `json:"account_id"`
	UserID    uint `json:"user_id"`
}

func (Membership) TableName() string {
	return "memberships"
}
