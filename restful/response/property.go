package response

import (
	"time"

	"github.com/solabsafrica/afrikanest/model"
)

type GetPropertyResposne struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreatePropertyResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUnitResponse struct {
	ID         string    `json:"id"`
	PropertyID string    `json:"proprty_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetPropertiesResponse struct {
	Pagination Pagination       `json:"pagination"`
	Properties []model.Property `json:"properties"`
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
