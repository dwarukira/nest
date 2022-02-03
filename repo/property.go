package repo

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/model"
	"gorm.io/gorm/clause"
)

type PropertyRepoWithContext func(ctx context.Context) PropertyRepo

type PropertyQuery struct {
	Name    *string
	OwnerID uuid.UUID
	Query   *string
	Offset  int
	Limit   int
}

type PropertyRepo interface {
	Create(model.Property) (model.Property, error)
	CreateUnit(model.Unit) (model.Unit, error)
	GetById(uuid.UUID) (model.Property, error)
	GetUnitById(uuid.UUID) (model.Unit, error)
	QueryProperties(query PropertyQuery) (properties []model.Property, total int64, err error)
	QueryUnits(query PropertyQuery) (units []model.Unit, total int64, err error)
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
	err := repo.db(repo.ctx).Preload("Units").Preload("Units.Leases").Preload("Units.Leases.Tenants").Preload("Units.Property").First(&property, "id = ?", id).Error()
	for _, unit := range *property.Units {
		unit.GetCurrentLease()
	}
	return property, err
}

func (repo *propertyRepoImpl) GetUnitById(id uuid.UUID) (model.Unit, error) {
	var property model.Unit
	err := repo.db(repo.ctx).Joins("Property").Preload("Leases").Preload("Leases.Tenants").First(&property, "units.id = ?", id).Error()
	return property, err
}

func (repo *propertyRepoImpl) QueryProperties(query PropertyQuery) ([]model.Property, int64, error) {
	var count int64
	properties := []model.Property{}
	db := repo.db(repo.ctx).Model(&model.Property{}).
		Offset(query.Offset).
		Limit(query.Limit).
		Where("owner_id = ?", query.OwnerID).Count(&count)

	if *query.Query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+*query.Query+"%", "%"+*query.Query+"%")
	}

	err := db.Count(&count).Find(&properties).Error()
	return properties, count, err
}

func (repo *propertyRepoImpl) QueryUnits(query PropertyQuery) ([]model.Unit, int64, error) {
	var count int64
	units := []model.Unit{}
	i, err := strconv.Atoi(*query.Query)
	if err != nil {
		// handle error
		i = 0
	}

	logger.Info(i)
	db := repo.db(repo.ctx).Debug().
		Joins("Property").Where("Property.owner_id", query.OwnerID).Preload("Leases.Unit").Preload("Leases.Tenants").Preload("Leases.LeaseCharge").Preload("Leases.LeaseCharge.LeaseChargesPayments")

	err = db.Model(&model.Unit{}).
		Where("units.name LIKE ? OR units.description LIKE ? OR \"Property\".\"name\" LIKE ? ", "%"+*query.Query+"%", "%"+*query.Query+"%", "%"+*query.Query+"%").
		Preload(clause.Associations).Preload("Leases.LeaseCharge.LeaseChargesPayments").
		Find(&units).Limit(query.Limit).Count(&count).Error()

	return units, count, err
}
