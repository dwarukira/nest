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

// swagger:model GetPropertiesResponse
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
	Status       []string     `json:"status"`
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewGetUnitsResponse(units []model.Unit, pagination Pagination) GetUnitsResponse {
	var unitsResponse []UnitsResponse
	for _, unit := range units {
		currentLease := unit.GetCurrentLease()
		// TODO: This is a hack to get the status of the unit.
		// consider moving this to the model.
		var status []string
		overdue := "https://badgen.net/badge/icon/overdue?label&labelColor=white&color=red"
		if currentLease != nil {
			for _, charge := range *currentLease.LeaseCharge {
				var leaseCharges int64
				for _, leaseCharge := range charge.LeaseChargesPayments {
					leaseCharges += leaseCharge.Amount
				}
				if charge.Amount-leaseCharges > 0 && charge.DueDate.Before(time.Now()) {
					if !contains(status, overdue) {
						status = append(status, overdue)
					}
				}
			}
		} else {
			status = append(status, "https://badgen.net/badge/icon/vacant?label&labelColor=white&color=yellow")
		}

		unitsResponse = append(unitsResponse, UnitsResponse{
			Unit:         unit,
			Status:       status,
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
