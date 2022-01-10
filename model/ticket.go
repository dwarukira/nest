package model

import "github.com/google/uuid"

type TicketStatus struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	Color       string `gorm:"column:color" json:"color"`
}

type IssueType struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}

type Ticket struct {
	Base
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	// Images      []string     `json:"images" gorm:"images"` # TODO: implement images
	Status      TicketStatus `json:"status"`
	StatusID    uuid.UUID    `gorm:"column:ticket_status_id" json:"status_id"`
	UnitID      uuid.UUID    `gorm:"column:unit_id" json:"unit_id"`
	Unit        Unit         `json:"unit"`
	IssueType   IssueType    `json:"issue_type"`
	IssueTypeID uuid.UUID    `gorm:"column:issue_type_id" json:"issue_type_id"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}

func (t *TicketStatus) TableName() string {
	return "ticket_statuses"
}
