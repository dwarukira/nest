package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/model"
)

type PropertyRepoWithContext func(ctx context.Context) PropertyRepo

type PropertyQuery struct {
	Name    *string
	OwnerID uuid.UUID

	Offset int
	Limit  int
}

type PropertyRepo interface {
	Create(model.Property) (model.Property, error)
	CreateUnit(model.Unit) (model.Unit, error)
	GetById(uuid.UUID) (model.Property, error)
	GetUnitById(uuid.UUID) (model.Unit, error)
	QueryProperties(query PropertyQuery) (properties []model.Property, total int64, err error)
}

type propertyRepoImpl struct {
	ctx context.Context
	db  db.DatabaseWithCtx
}

func NewPropertyRepoWithContext(db db.DatabaseWithCtx) PropertyRepoWithContext {
	return func(ctx context.Context) PropertyRepo {
		return &propertyRepoImpl{
			ctx: ctx,
			db:  db,
		}
	}
}

func (repo *propertyRepoImpl) Create(property model.Property) (model.Property, error) {
	err := repo.db(repo.ctx).Create(&property).Error()
	return property, err
}

func (repo *propertyRepoImpl) CreateUnit(unit model.Unit) (model.Unit, error) {
	err := repo.db(repo.ctx).Create(&unit).Error()
	return unit, err
}

func (repo *propertyRepoImpl) GetById(id uuid.UUID) (model.Property, error) {
	var property model.Property
	err := repo.db(repo.ctx).First(&property, "id = ?", id).Error()
	return property, err
}

func (repo *propertyRepoImpl) GetUnitById(id uuid.UUID) (model.Unit, error) {
	var property model.Unit
	err := repo.db(repo.ctx).Joins("Property").First(&property, "units.id = ?", id).Error()
	return property, err
}

func (repo *propertyRepoImpl) QueryProperties(query PropertyQuery) ([]model.Property, int64, error) {
	var count int64
	properties := []model.Property{}
	db := repo.db(repo.ctx).Model(&model.Property{}).
		Offset(query.Offset).
		Limit(query.Limit).
		Where("owner_id = ?", query.OwnerID).Count(&count)
	err := db.Find(&properties).Error()
	return properties, count, err
}
