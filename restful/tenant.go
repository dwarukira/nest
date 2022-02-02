package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/solabsafrica/afrikanest/restful/middlewares"
	"github.com/solabsafrica/afrikanest/restful/response"
	"github.com/solabsafrica/afrikanest/service"
)

type tenantController struct {
	tenantService service.TenantServiceWithContext
}

func NewTenantController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, tenantService service.TenantServiceWithContext) {
	controller := &tenantController{
		tenantService: tenantService,
	}

	v1 := group.Group("/v1")
	v1.GET("/tenants", authChecker.Check, controller.GetTenants)
}

// swagger:route GET /tenants Tenants getTenants
//
// Get all tenants
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http, https
//
// Responses:
//  200: TenantsResponse
//  401: ErrorResponse
//  500: ErrorResponse

func (controller *tenantController) GetTenants(ctx *gin.Context) {
	indentity, err := GetIndentityFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	tenants, err := controller.tenantService(ctx).GetTenantsForUser(indentity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, response.NewTenantsResponse(tenants, response.NewPagination(0, 0, 0)))
}
