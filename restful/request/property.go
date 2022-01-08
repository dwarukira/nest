package request

import (
	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/model"
)

type CreatePropertyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (createPropertyRequest CreatePropertyRequest) Validate() error {
	if len(createPropertyRequest.Name) == 0 {
		return exceptions.PropertyCreateFailed.SetMessage("name must be provided")
	}
	return nil
}

func (createPropertyRequest CreatePropertyRequest) ToProperty() (model.Property, error) {
	if err := createPropertyRequest.Validate(); err != nil {
		return model.Property{}, err
	}
	return model.Property{
		Name:        createPropertyRequest.Name,
		Description: createPropertyRequest.Description,
	}, nil
}

type CreateUnitRequest struct {
	CreatePropertyRequest
	PropertyID string `json:"property_id"`
}

func (createUnitRequest CreateUnitRequest) Validate() error {
	if len(createUnitRequest.Name) == 0 {
		return exceptions.UnitCreateFailed.SetMessage("name must be provided")
	}

	if len(createUnitRequest.PropertyID) == 0 {
		return exceptions.UnitCreateFailed.SetMessage("property_id must be provided")
	}

	return nil
}

func (createUnitRequest CreateUnitRequest) ToUnit() (model.Unit, error) {
	if err := createUnitRequest.Validate(); err != nil {
		return model.Unit{}, err
	}

	i, err := uuid.Parse(createUnitRequest.PropertyID)
	if err != nil {
		return model.Unit{}, exceptions.UnitCreateFailed.SetMessage("property_id must be valid")
	}
	return model.Unit{
		Name:        createUnitRequest.Name,
		Description: createUnitRequest.Description,
		PropertyID:  i,
	}, nil

}
