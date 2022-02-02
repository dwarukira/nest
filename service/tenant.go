package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/repo"
)

type TenantServiceWithContext func(ctx context.Context) TenantService

type TenantService interface {
	Create(tenant model.Tenant) (model.Tenant, error)
	GetById(id uuid.UUID) (model.Tenant, error)
	GetTenantByOwnerId(ownerId uuid.UUID) (model.Tenant, error)
	GetTenantsForUser(userId uuid.UUID) ([]model.Tenant, error)
}

type tenantServiceImpl struct {
	ctx        context.Context
	tenantRepo repo.TenantRepoWithContext
	userRepo   repo.UserRepoWithContext
}

func (tenantService *tenantServiceImpl) Create(tenant model.Tenant) (model.Tenant, error) {
	tenant, err := tenantService.tenantRepo(tenantService.ctx).Create(tenant)
	if err != nil {
		return tenant, exceptions.TenantCreateFaild.Wrap(err).SetMessage(err.Error())
	}
	return tenant, err
}

func (tenantService *tenantServiceImpl) GetById(id uuid.UUID) (model.Tenant, error) {
	tenant, err := tenantService.tenantRepo(tenantService.ctx).GetTenantById(id)
	if err != nil {
		return tenant, exceptions.TenantNotExists.Wrap(err).SetMessage(err.Error())
	}
	return tenant, err
}

func (tenantService *tenantServiceImpl) GetTenantByOwnerId(ownerId uuid.UUID) (model.Tenant, error) {
	tenant, err := tenantService.tenantRepo(tenantService.ctx).GetTenantByOwnerId(ownerId)
	if err != nil {
		return tenant, exceptions.TenantNotExists.Wrap(err).SetMessage(err.Error())
	}
	return tenant, err
}

func (tenantService *tenantServiceImpl) GetTenantsForUser(userId uuid.UUID) ([]model.Tenant, error) {
	// tenant, err := tenantService.tenantRepo(tenantService.ctx).GetTenantForUser(userId)
	// if err != nil {
	// 	return tenant, exceptions.TenantNotExists.Wrap(err).SetMessage(err.Error())
	// }
	// return tenant, err
	user, err := tenantService.userRepo(tenantService.ctx).GetById(userId)
	if err != nil {
		return nil, err
	}

	var tenants []model.Tenant

	for _, account := range user.Accounts {
		logger.Info("account", account)
		for _, tenant := range account.Tenants {
			tenants = append(tenants, *tenant)
		}
	}

	return tenants, nil
}

func NewTenantServiceWithContext(tenantRepo repo.TenantRepoWithContext, user repo.UserRepoWithContext) TenantServiceWithContext {
	return func(ctx context.Context) TenantService {
		return &tenantServiceImpl{
			ctx:        ctx,
			tenantRepo: tenantRepo,
			userRepo:   user,
		}
	}
}
