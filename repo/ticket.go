package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/model"
)

type TicketRepoWithContext func(ctx context.Context) TicketRepo

type TicketQuery struct {
	Name *string

	Offset int
	Limit  int
}

type TicketRepo interface {
	Create(model.Ticket) (model.Ticket, error)
	GetById(uuid.UUID) (model.Ticket, error)
	Save(model.Ticket) error
	QueryTickets(query TicketQuery) (tickets []model.Ticket, total int64, err error)
}

type ticketRepoImpl struct {
	ctx context.Context
	db  db.DatabaseWithCtx
}

func (repo *ticketRepoImpl) Create(ticket model.Ticket) (model.Ticket, error) {
	err := repo.db(repo.ctx).Create(&ticket).Error()
	return ticket, err
}

func (repo *ticketRepoImpl) Save(ticket model.Ticket) error {
	return repo.db(repo.ctx).Save(&ticket).Error()
}

func (repo *ticketRepoImpl) GetById(id uuid.UUID) (model.Ticket, error) {
	var ticket model.Ticket
	err := repo.db(repo.ctx).First(&ticket, "id = ?", id).Error()
	return ticket, err
}

func (repo *ticketRepoImpl) QueryTickets(query TicketQuery) (tickets []model.Ticket, total int64, err error) {
	err = repo.db(repo.ctx).Model(&model.Ticket{}).Count(&total).Error()
	if err != nil {
		return nil, 0, err
	}
	err = repo.db(repo.ctx).Offset(query.Offset).Limit(query.Limit).Find(&tickets).Error()
	return tickets, total, err
}

func NewTicketRepoWithContext(db db.DatabaseWithCtx) TicketRepoWithContext {
	return func(ctx context.Context) TicketRepo {
		return &ticketRepoImpl{
			ctx: ctx,
			db:  db,
		}
	}
}
