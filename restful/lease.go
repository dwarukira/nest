package restful

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/restful/middlewares"
	"github.com/solabsafrica/afrikanest/restful/request"
	"github.com/solabsafrica/afrikanest/restful/response"
	"github.com/solabsafrica/afrikanest/service"
)

type leaseController struct {
	leaseService    service.LeaseServiceWithContext
	propertyService service.PropertyServiceWithContext
}

func NewLeaseController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, leaseService service.LeaseServiceWithContext, propertyService service.PropertyServiceWithContext) {
	controller := &leaseController{
		leaseService:    leaseService,
		propertyService: propertyService,
	}

	v1 := group.Group("/v1")

	v1.POST("/leases", authChecker.Check, controller.CreateLeaseHandler)
	v1.GET("/leases/:id", authChecker.Check, controller.GetLeaseHandler)

}

func (ctrl *leaseController) CreateLeaseHandler(ctx *gin.Context) {
	var createLeaseRequest request.CreateLeaseRequest
	if err := ctx.ShouldBindJSON(&createLeaseRequest); err != nil {
		logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	lease, err := createLeaseRequest.ToLease()
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	lease.LeaseStatus = model.DRAFT
	random_lease_number := uuid.New().String()
	lease.LeaseNumber = random_lease_number
	// ensure unit exist and we own it
	unit, err := ctrl.propertyService(ctx).GetUnitById(lease.UnitID)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	logger.Info(unit.Property.Name)
	identity, _ := GetIndentityFromContext(ctx)

	if unit.Property.OwnerID != identity {
		logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if lease, err = ctrl.leaseService(ctx).Create(lease); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusCreated, response.NewCreateLeaseResponse(lease))
}

func (ctrl *leaseController) GetLeaseHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("id")))
		return
	}

	lease, err := ctrl.leaseService(ctx).GetLeaseById(id)
	if err != nil {
		if errors.Is(err, exceptions.LeaseNotExists) {
			ctx.JSON(http.StatusNotFound, err)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	// TODO: check access when we have tenants
	ctx.JSON(http.StatusOK, lease)

}
