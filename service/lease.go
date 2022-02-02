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
	GetLeaseCharge(uuid.UUID) (model.LeaseCharge, error)
	CreateLeaseCharge(model.LeaseCharge) (model.LeaseCharge, error)
	CreateLeaseChargePayment(model.LeaseChargePayment) (model.LeaseChargePayment, error)
	// GetLeaseBalance(uuid.UUID) (float64, error)
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
	var s []model.LeaseCharge
	for _, e := range *lease.LeaseCharge {
		var leaseCharges int64
		for _, p := range e.LeaseChargesPayments {
			leaseCharges += p.Amount
		}
		e.RecivedAmount = leaseCharges
		e.RemainingAmount = e.Amount - leaseCharges
		s = append(s, e)
	}
	lease.LeaseCharge = &s
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return lease, exceptions.LeaseNotExists.Wrap(err).SetMessage(err.Error())
	}
	return lease, err
}

func (leaseService *leaseServiceImpl) GetLeasesForUnit(page int, pageSize int, unitId uuid.UUID) ([]model.Lease, int64, error) {
	query := repo.LeaseQuery{Offset: page * pageSize, Limit: pageSize}
	return leaseService.leaseRepo(leaseService.ctx).QueryLeases(query)
}

func (leaseService *leaseServiceImpl) GetLeaseCharge(leaseChargeId uuid.UUID) (model.LeaseCharge, error) {
	return leaseService.leaseRepo(leaseService.ctx).GetLeaseCharge(leaseChargeId)
}

func (leaseService *leaseServiceImpl) CreateLeaseChargePayment(leaseChargePayment model.LeaseChargePayment) (model.LeaseChargePayment, error) {
	leaseCharge, err := leaseService.GetLeaseCharge(leaseChargePayment.LeaseChargeID)
	if err != nil {
		return model.LeaseChargePayment{}, err
	}

	if leaseCharge.RemainingAmount < leaseChargePayment.Amount {
		return model.LeaseChargePayment{}, exceptions.LeaseChargePaymentAmountExceedsRemainingAmount.Wrap(err).SetMessage("Lease charge payment amount exceeds remaining amount")
	}

	leaseChargePayment, err = leaseService.leaseRepo(leaseService.ctx).CreateLeaseChargePayment(leaseChargePayment)
	if err != nil {
		return leaseChargePayment, exceptions.LeaseChargePaymentCreateFaild.Wrap(err).SetMessage(err.Error())
	}
	return leaseChargePayment, nil
}

func (leaseService *leaseServiceImpl) GetLeaseBalance(leaseId uuid.UUID) (int64, error) {
	v, err := leaseService.leaseRepo(leaseService.ctx).GetLeaseCharge(leaseId)
	if err != nil {
		return 0, err
	}

	return v.RemainingAmount, nil
}

func (leaseService *leaseServiceImpl) CreateLeaseCharge(leaseCharge model.LeaseCharge) (model.LeaseCharge, error) {
	leaseCharge, err := leaseService.leaseRepo(leaseService.ctx).CreateLeaseCharge(leaseCharge)
	if err != nil {
		return leaseCharge, exceptions.LeaseChargeCreateFaild.Wrap(err).SetMessage(err.Error())
	}
	return leaseCharge, nil
}
