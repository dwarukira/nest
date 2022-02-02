package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/model"
)

type LeaseRepoWithContext func(ctx context.Context) LeaseRepo

type LeaseQuery struct {
	UnitID uuid.UUID

	Offset int
	Limit  int
}

type LeaseRepo interface {
	Create(model.Lease) (model.Lease, error)
	GetLeaseById(uuid.UUID) (model.Lease, error)
	QueryLeases(query LeaseQuery) (lease []model.Lease, total int64, err error)
	GetLeaseCharge(uuid.UUID) (model.LeaseCharge, error)
	CreateLeaseCharge(model.LeaseCharge) (model.LeaseCharge, error)
	CreateLeaseChargePayment(model.LeaseChargePayment) (model.LeaseChargePayment, error)
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

func (repo *leaseRepoImpl) GetLeaseById(id uuid.UUID) (model.Lease, error) {
	var lease model.Lease
	err := repo.db(repo.ctx).Debug().Preload("Unit").
		Preload("Unit.Property").
		Preload("Tenants").
		Preload("LeaseCharge").
		Preload("LeaseCharge.LeaseChargesPayments").
		Preload("Tenants.Lease").
		Preload("Tenants.Lease.Unit.Property").
		First(&lease, "id = ?", id).Error()
	return lease, err
}

func (repo *leaseRepoImpl) GetLeaseCharge(id uuid.UUID) (model.LeaseCharge, error) {
	var leaseCharge model.LeaseCharge

	err := repo.db(repo.ctx).Debug().
		Preload("Lease").
		Preload("LeaseChargesPayments").
		First(&leaseCharge, "id = ?", id).Error()
	var leaseCharges int64
	for _, leaseCharge := range leaseCharge.LeaseChargesPayments {
		leaseCharges += leaseCharge.Amount
	}
	leaseCharge.RemainingAmount = leaseCharge.Amount - leaseCharges
	leaseCharge.RecivedAmount = leaseCharges
	return leaseCharge, err
}

func (repo *leaseRepoImpl) QueryLeases(query LeaseQuery) ([]model.Lease, int64, error) {
	var count int64
	leases := []model.Lease{}
	db := repo.db(repo.ctx).Debug().Model(&model.Lease{}).
		Offset(query.Offset).
		Limit(query.Limit).
		Where("unit_id = ?", query.UnitID).
		Count(&count)
	err := db.Find(&leases).Error()
	return leases, count, err
}

func (repo *leaseRepoImpl) CreateLeaseChargePayment(leaseChargePayment model.LeaseChargePayment) (model.LeaseChargePayment, error) {
	err := repo.db(repo.ctx).Create(&leaseChargePayment).Error()
	return leaseChargePayment, err
}

func (repo *leaseRepoImpl) CreateLeaseCharge(leaseCharge model.LeaseCharge) (model.LeaseCharge, error) {
	err := repo.db(repo.ctx).Create(&leaseCharge).Error()
	return leaseCharge, err
}
