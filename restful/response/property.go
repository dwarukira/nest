package response

import (
	"time"

	"github.com/solabsafrica/afrikanest/model"
)

type GetPropertyResposne struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// swagger:model CreatePropertyResponse
type CreatePropertyResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// swagger:model CreateUnitResponse
type CreateUnitResponse struct {
	ID         string    `json:"id"`
	PropertyID string    `json:"proprty_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetPropertiesResponse struct {
	Pagination Pagination       `json:"pagination"`
	Properties []model.Property `json:"properties"`
}

// swagger:model GetPropertyResponse
type GetPropertyResponse struct {
	Property model.Property `json:"property"`
}

// swagger:model GetUnitsResponse
type GetUnitsResponse struct {
	Pagination Pagination      `json:"pagination"`
	Units      []UnitsResponse `json:"units"`
}

// swagger:model UnitsResponse
type UnitsResponse struct {
	model.Unit
	CurrentLease *model.Lease `json:"current_lease"`
}

func NewUnitResponse(unit model.Unit) UnitsResponse {
	currentLease := unit.GetCurrentLease()

	return UnitsResponse{
		Unit:         unit,
		CurrentLease: currentLease,
	}
}

type GetUnitResponse struct {
}

func NewCreatePropertyResponse(property model.Property) CreatePropertyResponse {
	return CreatePropertyResponse{
		ID:        property.ID.String(),
		CreatedAt: property.CreatedAt,
	}
}

func NewCreateUnitResponse(unit model.Unit) CreateUnitResponse {
	return CreateUnitResponse{
		ID:         unit.ID.String(),
		PropertyID: unit.PropertyID.String(),
		CreatedAt:  unit.CreatedAt,
	}
}

func NewGetPropertiesResponse(properties []model.Property, pagination Pagination) GetPropertiesResponse {
	return GetPropertiesResponse{
		Properties: properties,
		Pagination: pagination,
	}
}

func NewGetUnitsResponse(units []model.Unit, pagination Pagination) GetUnitsResponse {
	var unitsResponse []UnitsResponse
	for _, unit := range units {
		currentLease := unit.GetCurrentLease()

		unitsResponse = append(unitsResponse, UnitsResponse{
			Unit:         unit,
			CurrentLease: currentLease,
		})
	}
	return GetUnitsResponse{
		Units:      unitsResponse,
		Pagination: pagination,
	}
}

func NewGetPropertyResponse(property model.Property) GetPropertyResponse {
	return GetPropertyResponse{
		Property: property,
	}
}
