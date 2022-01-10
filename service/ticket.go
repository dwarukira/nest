package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/repo"
)

type TicketServiceWithContext func(ctx context.Context) TicketService

type TicketService interface {
	Create(ticket model.Ticket) (model.Ticket, error)
	GetById(id uuid.UUID) (model.Ticket, error)
}

type ticketServiceImpl struct {
	ctx         context.Context
	tickectRepo repo.TicketRepoWithContext
}

func NewTicketServiceWithContext(tickectRepo repo.TicketRepoWithContext) TicketServiceWithContext {
	return func(ctx context.Context) TicketService {
		return &ticketServiceImpl{
			ctx:         ctx,
			tickectRepo: tickectRepo,
		}
	}
}

func (ticketService *ticketServiceImpl) Create(ticket model.Ticket) (model.Ticket, error) {
	ticket, err := ticketService.tickectRepo(ticketService.ctx).Create(ticket)
	if err != nil {
		return ticket, exceptions.TicketCreateFaild.Wrap(err).SetMessage(err.Error())
	}
	return ticket, err
}

func (ticketService *ticketServiceImpl) GetById(id uuid.UUID) (model.Ticket, error) {
	ticket, err := ticketService.tickectRepo(ticketService.ctx).GetById(id)
	if err != nil {
		return ticket, exceptions.TicketNotExists.Wrap(err).SetMessage(err.Error())
	}
	return ticket, err
}
