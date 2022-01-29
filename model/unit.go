package model

import "github.com/google/uuid"

// swagger:model Unit
type Unit struct {
	Base
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	DefaultRent int       `gorm:"column:default_rent" json:"default_rent"`
	Property    Property  `json:"property"`
	PropertyID  uuid.UUID `gorm:"column:property_id" json:"property_id"`
	Leases      []Lease   `json:"leases"`
}

func (u Unit) TableName() string {
	return "units"
}

func (u Unit) GetCurrentLease() *Lease {
	var currentLease *Lease

	for _, v := range u.Leases {
		if v.LeaseStatus == "ACTIVE" {
			currentLease = &v
		}
	}

	if currentLease == nil {
		return nil
	}

	return currentLease
}
