package response

import "github.com/solabsafrica/afrikanest/model"

// swagger:model TenantsResponse
type TenantsResponse struct {
	Tenants    []model.Tenant `json:"tenants"`
	Pagination Pagination     `json:"pagination"`
}

func NewTenantsResponse(tenant []model.Tenant, pagination Pagination) TenantsResponse {
	return TenantsResponse{
		Tenants:    tenant,
		Pagination: pagination,
	}
}
