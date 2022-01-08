package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/repo"
	"gorm.io/gorm"
)

type LeaseServiceWithContext func(ctx context.Context) LeaseService

type LeaseService interface {
	Create(model.Lease) (model.Lease, error)
	GetLeaseById(uuid.UUID) (model.Lease, error)
	GetLeasesForUnit(page int, pageSize int, unitId uuid.UUID) ([]model.Lease, int64, error)
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

func (leaseService *leaseServiceImpl) GetLeaseById(id uuid.UUID) (model.Lease, error) {
	lease, err := leaseService.leaseRepo(leaseService.ctx).GetLeaseById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lease, exceptions.LeaseNotExists.Wrap(err).SetMessage(err.Error())
	}
	return lease, err
}

func (leaseService *leaseServiceImpl) GetLeasesForUnit(page int, pageSize int, unitId uuid.UUID) ([]model.Lease, int64, error) {
	query := repo.LeaseQuery{Offset: page * pageSize, Limit: pageSize}
	return leaseService.leaseRepo(leaseService.ctx).QueryLeases(query)
}
