package response

import (
	"time"

	"github.com/solabsafrica/afrikanest/model"
)

type CreateTicketResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCreateTicketResponse(ticket model.Ticket) CreateTicketResponse {
	return CreateTicketResponse{
		ID:        ticket.ID.String(),
		CreatedAt: ticket.CreatedAt,
	}
}
