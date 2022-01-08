package service

import (
	"context"

	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/repo"
)

type LeaseServiceWithContext func(ctx context.Context) LeaseService

type LeaseService interface {
	Create(model.Lease) (model.Lease, error)
}

type leaseServiceImpl struct {
	ctx       context.Context
	leaseRepo repo.LeaseRepoWithContext
}

func NewLeaseServiceWithContext(leaseRepo repo.LeaseRepoWithContext) LeaseServiceWithContext {
	return func(ctx context.Context) LeaseService {
		return &leaseServiceImpl{
			ctx:       ctx,
			leaseRepo: leaseRepo,
		}
	}
}

func (leaseService *leaseServiceImpl) Create(lease model.Lease) (model.Lease, error) {
	lease, err := leaseService.leaseRepo(leaseService.ctx).Create(lease)
	if err != nil {
		return lease, exceptions.LeaseCreateFaild.Wrap(err).SetMessage(err.Error())
	}
	return lease, nil
}
