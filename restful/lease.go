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
	smsService      service.SmsServiceWithContext
}

func NewLeaseController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, leaseService service.LeaseServiceWithContext, propertyService service.PropertyServiceWithContext, smsService service.SmsServiceWithContext) {
	controller := &leaseController{
		leaseService:    leaseService,
		propertyService: propertyService,
		smsService:      smsService,
	}

	v1 := group.Group("/v1")

	v1.POST("/leases", authChecker.Check, controller.CreateLeaseHandler)
	v1.GET("/leases/:id", authChecker.Check, controller.GetLeaseHandler)
	v1.GET("/leases/:id/tenants", authChecker.Check, controller.GetLeaseTenantsHandler)
}

// swagger:route POST /leases lease createLease
//
//  Create a lease
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	201: CreateLeaseResponse
//  401: ErrorResponse
//  500: ErrorResponse
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

	lease.LeaseStatus = model.ACTIVE
	random_lease_number := uuid.New().String()
	lease.LeaseNumber = random_lease_number

	// ensure unit exist and we own it
	unit, err := ctrl.propertyService(ctx).GetUnitById(lease.UnitID)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	identity, _ := GetIndentityFromContext(ctx)
	if unit.Property.OwnerID != identity {
		logger.Error(errors.New("you are not the owner of this property"))
		ctx.JSON(http.StatusUnauthorized, exceptions.LeaseCreateFaild.SetMessage("you are not the owner of this property"))
		return
	}
	for _, l := range unit.Leases {
		if l.LeaseStatus == model.ACTIVE {
			logger.Error(errors.New("unit is already leased"))
			ctx.JSON(http.StatusUnauthorized, exceptions.LeaseCreateFaild.SetMessage("unit is already leased, consider cancelling the existing lease"))
			return
		}
	}
	if lease, err = ctrl.leaseService(ctx).Create(lease); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
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
	// err = ctrl.smsService(ctx).Send("+254728165763", "Hello World")

	ctx.JSON(http.StatusOK, lease)

}

func (ctrl *leaseController) GetLeaseTenantsHandler(ctx *gin.Context) {

}
