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

type PropertyServiceWithContext func(ctx context.Context) PropertyService

type PropertyService interface {
	Create(model.Property) (model.Property, error)
	CreateUnit(model.Unit) (model.Unit, error)
	GetPropertyById(uuid.UUID) (model.Property, error)
	ListUserProperties(page int, pageSize int, userId uuid.UUID) ([]model.Property, int64, error)
	ListUserUnits(page int, pageSize int, userId uuid.UUID, query string) ([]model.Unit, int64, error)
	GetUnitById(uuid.UUID) (model.Unit, error)
}

type propertyServiceImpl struct {
	ctx          context.Context
	propertyRepo repo.PropertyRepoWithContext
}

func (propertyService *propertyServiceImpl) Create(property model.Property) (model.Property, error) {
	property, err := propertyService.propertyRepo(propertyService.ctx).Create(property)
	if err != nil {
		return property, exceptions.PropertyCreateFailed.Wrap(err).SetMessage(err.Error())
	}

	return property, nil

}

func (propertyService *propertyServiceImpl) CreateUnit(unit model.Unit) (model.Unit, error) {
	unit, err := propertyService.propertyRepo(propertyService.ctx).CreateUnit(unit)
	if err != nil {
		return unit, exceptions.UnitCreateFailed.Wrap(err).SetMessage(err.Error())
	}

	return unit, nil

}

func (propertyService *propertyServiceImpl) GetPropertyById(id uuid.UUID) (model.Property, error) {
	property, err := propertyService.propertyRepo(propertyService.ctx).GetById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return property, exceptions.PropertyNotExists.Wrap(err)
	}
	return property, err
}

func (propertyService *propertyServiceImpl) GetUnitById(id uuid.UUID) (model.Unit, error) {
	unit, err := propertyService.propertyRepo(propertyService.ctx).GetUnitById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return unit, exceptions.UnitNotExists.Wrap((err))
	}
	return unit, err
}

func (propertyService *propertyServiceImpl) ListUserProperties(page int, pageSize int, userID uuid.UUID) ([]model.Property, int64, error) {
	query := repo.PropertyQuery{Offset: page * pageSize, Limit: pageSize, OwnerID: userID}
	return propertyService.propertyRepo(propertyService.ctx).QueryProperties(query)
}

func (propertyService *propertyServiceImpl) ListUserUnits(page int, pageSize int, userID uuid.UUID, queryParam string) ([]model.Unit, int64, error) {
	query := repo.PropertyQuery{Offset: page * pageSize, Limit: pageSize, OwnerID: userID, Query: &queryParam}
	return propertyService.propertyRepo(propertyService.ctx).QueryUnits(query)
}

func NewPropertyServiceWithContext(propertyRepo repo.PropertyRepoWithContext) PropertyServiceWithContext {
	return func(ctx context.Context) PropertyService {
		return &propertyServiceImpl{
			ctx:          ctx,
			propertyRepo: propertyRepo,
		}
	}
}
