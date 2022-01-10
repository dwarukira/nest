package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/model"
)

type TenantRepoWithContext func(ctx context.Context) TenantRepo

type TenantQuery struct {
	LeaseID *uuid.UUID

	Offset int
	Limit  int
}

type TenantRepo interface {
	Create(model.Tenant) (model.Tenant, error)
	GetTenantById(uuid.UUID) (model.Tenant, error)
	QueryTenants(query TenantQuery) (tenants []model.Tenant, total int64, err error)
}

type tenantRepoImpl struct {
	ctx context.Context
	db  db.DatabaseWithCtx
}

func (repo *tenantRepoImpl) Create(tenant model.Tenant) (model.Tenant, error) {
	err := repo.db(repo.ctx).Create(&tenant).Error()
	return tenant, err
}

func (repo *tenantRepoImpl) GetTenantById(id uuid.UUID) (model.Tenant, error) {
	var tenant model.Tenant
	err := repo.db(repo.ctx).Preload("Lease").Preload("Lease.Unit").First(&tenant, "id = ?", id).Error()
	return tenant, err
}

func (repo *tenantRepoImpl) QueryTenants(query TenantQuery) ([]model.Tenant, int64, error) {
	var count int64
	tenants := []model.Tenant{}
	db := repo.db(repo.ctx).Model(&model.Tenant{}).
		Offset(query.Offset).
		Limit(query.Limit).
		Where("lease_id = ?", query.LeaseID).
		Count(&count)
	err := db.Find(&tenants).Error()
	return tenants, count, err
}

func NewTenantRepoWithContext(db db.DatabaseWithCtx) TenantRepoWithContext {
	return func(ctx context.Context) TenantRepo {
		return &tenantRepoImpl{
			ctx: ctx,
			db:  db,
		}
	}
}
