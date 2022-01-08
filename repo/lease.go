package repo

import (
	"context"

	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/model"
)

type LeaseRepoWithContext func(ctx context.Context) LeaseRepo

type LeaseRepo interface {
	Create(model.Lease) (model.Lease, error)
}

type leaseRepoImpl struct {
	ctx context.Context
	db  db.DatabaseWithCtx
}

func NewLeaseRepoWithContext(db db.DatabaseWithCtx) LeaseRepoWithContext {
	return func(ctx context.Context) LeaseRepo {
		return &leaseRepoImpl{
			ctx: ctx,
			db:  db,
		}
	}
}

func (repo *leaseRepoImpl) Create(lease model.Lease) (model.Lease, error) {
	err := repo.db(repo.ctx).Create(&lease).Error()
	return lease, err
}
