package request

import (
	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
)

type CreateTicketRequest struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	UnitID            string `json:"unit_id"`
	TicketStatusID    string `json:"ticket_status_id"`
	TicketIssueTypeID string `json:"ticket_issue_type_id"`
}

func (createTicketRequest CreateTicketRequest) Validate() error {
	if len(createTicketRequest.Title) == 0 {
		return exceptions.TicketCreateFaild.SetMessage("title must be provided")
	}

	if len(createTicketRequest.UnitID) == 0 {
		return exceptions.TicketCreateFaild.SetMessage("unit_id must be provided")
	}

	return nil
}

func (createTicketRequest CreateTicketRequest) ToTicket() (model.Ticket, error) {
	if err := createTicketRequest.Validate(); err != nil {
		return model.Ticket{}, err
	}

	i, err := uuid.Parse(createTicketRequest.UnitID)
	if err != nil {
		return model.Ticket{}, exceptions.TicketCreateFaild.SetMessage("unit_id must be valid")
	}

	statusID, err := uuid.Parse(createTicketRequest.TicketStatusID)
	if err != nil {
		return model.Ticket{}, exceptions.TicketCreateFaild.SetMessage("ticket_type_id must be valid")
	}
	issueTypeID, err := uuid.Parse(createTicketRequest.TicketIssueTypeID)
	if err != nil {
		return model.Ticket{}, exceptions.TicketCreateFaild.SetMessage("ticket_issue_type_id must be valid")
	}

	return model.Ticket{
		Title:       createTicketRequest.Title,
		Description: createTicketRequest.Description,
		StatusID:    statusID,
		UnitID:      i,
		IssueTypeID: issueTypeID,
	}, nil
}
