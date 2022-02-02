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
	emailService    service.EmailServiceWithContext
}

func NewLeaseController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, leaseService service.LeaseServiceWithContext, propertyService service.PropertyServiceWithContext, smsService service.SmsServiceWithContext, emailService service.EmailServiceWithContext) {
	controller := &leaseController{
		leaseService:    leaseService,
		propertyService: propertyService,
		smsService:      smsService,
		emailService:    emailService,
	}

	v1 := group.Group("/v1")

	v1.POST("/leases", authChecker.Check, controller.CreateLeaseHandler)
	v1.GET("/leases/:id", authChecker.Check, controller.GetLeaseHandler)
	v1.GET("/leases/:id/tenants", authChecker.Check, controller.GetLeaseTenantsHandler)
	v1.GET("/leases/:id/charges", authChecker.Check, controller.GetLeaseCharges)
	v1.GET("/leases/:id/charges/:leaseChargeId", authChecker.Check, controller.GetLeaseCharge)
	v1.POST("/leases/:id/charges", authChecker.Check, controller.CreateLeaseCharge)
	v1.POST("/leases/:id/charges/:leaseChargeId/payments", authChecker.Check, controller.CreateLeaseChargePayment)
	v1.GET("/leases/:id/balance", authChecker.Check, controller.GetLeaseBalance)
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
	// TODO move this to a background job
	// send sms to tenant
	for _, tenant := range *lease.Tenants {
		logger.Info("sending sms to tenant", tenant.PhoneNumber)
		if tenant.PhoneNumber != "" {
			message := "Hi " + tenant.FirstName + ", your lease is now active. download the app to view your lease details"
			err = ctrl.smsService(ctx).Send(tenant.PhoneNumber, message)
			if err != nil {
				logger.Error(err)
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
		}

		if tenant.Email != "" {
			err = ctrl.emailService(ctx).SendTenantWelcomeEmail(tenant)
			if err != nil {
				logger.Error(err)
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
		}
	}

	// generate a new deposit charge for the tenant
	depositCharge := model.LeaseCharge{
		LeaseID:    lease.ID,
		Amount:     int64(lease.SecurityDeposit),
		ChargeType: model.SECURITY_DEPOSIT,
		DueDate:    lease.StartDate,
	}

	if _, err := ctrl.leaseService(ctx).CreateLeaseCharge(depositCharge); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	//  generate a new rent charge for the tenant
	rentCharge := model.LeaseCharge{
		LeaseID:    lease.ID,
		Amount:     int64(lease.MonthlyRent),
		ChargeType: model.RENT,
		DueDate:    lease.StartDate,
	}

	if _, err := ctrl.leaseService(ctx).CreateLeaseCharge(rentCharge); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewCreateLeaseResponse(lease))
}

// swagger:route GET /leases/{id} lease getLease
//
// Get a lease
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	201: LeaseResponse
//  401: ErrorResponse
//  500: ErrorResponse
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

func (ctrl *leaseController) GetLeaseTenantsHandler(ctx *gin.Context) {

}

func (ctrl *leaseController) GetLeaseCharges(ctx *gin.Context) {
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

	logger.Debug(lease)
}

// swagger:route GET /leases/{id}/charges/{leaseChargeId} lease getLeaseCharge
//
// Get a lease charges
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	200: LeaseChargeResponse
//  401: ErrorResponse
//  500: ErrorResponse

func (ctrl *leaseController) GetLeaseCharge(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("leaseChargeId"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("leaseChargeId")))
		return
	}
	leaseCharge, err := ctrl.leaseService(ctx).GetLeaseCharge(id)
	if err != nil {
		if errors.Is(err, exceptions.LeaseNotExists) {
			ctx.JSON(http.StatusNotFound, err)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response.NewLeaseChargeResponse(leaseCharge))
}

// swagger:route POST /leases/{id}/charges lease createLeaseCharge
//
//  Create a lease charge
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	201: LeaseChargeResponse
//  401: ErrorResponse
//  500: ErrorResponse
func (ctrl *leaseController) CreateLeaseCharge(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("id")))
		return
	}

	var createLeaseChargeRequest request.CreateLeaseChargeRequestBody
	if err := ctx.ShouldBindJSON(&createLeaseChargeRequest); err != nil {
		logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	leaseCharge, err := createLeaseChargeRequest.ToModel(id)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	leaseCharge.LeaseID = id

	leaseCharge, err = ctrl.leaseService(ctx).CreateLeaseCharge(leaseCharge)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewLeaseChargeResponse(leaseCharge))
}

// swagger:route POST /leases/{id}/charges/{leaseChargeId}/payments lease createLeaseChargePayment
//
//  Create a lease charge
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	201: LeaseChargePaymentResponse
//  401: ErrorResponse
//  500: ErrorResponse
func (ctrl *leaseController) CreateLeaseChargePayment(ctx *gin.Context) {
	var createLeaseChargePaymentRequest request.CreateLeaseChargePaymentRequestBody
	id, err := uuid.Parse(ctx.Param("leaseChargeId"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("leaseChargeId")))
		return
	}
	if err := ctx.ShouldBindJSON(&createLeaseChargePaymentRequest); err != nil {
		logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	leaseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("leaseID")))
		return
	}
	lease, err := ctrl.leaseService(ctx).GetLeaseById(leaseID)
	if err != nil {
		if errors.Is(err, exceptions.LeaseNotExists) {
			ctx.JSON(http.StatusNotFound, err)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	leaseChargePayment, err := createLeaseChargePaymentRequest.ToModel(id)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	leaseChargePayment, err = ctrl.leaseService(ctx).CreateLeaseChargePayment(leaseChargePayment)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	for _, tenant := range *lease.Tenants {
		if tenant.Email != "" && createLeaseChargePaymentRequest.SentEmailToTenant {
			if err := ctrl.emailService(ctx).SendLeaseChargePaymentEmail(tenant, lease, leaseChargePayment.Amount); err != nil {
				logger.Error(err)
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
		}
	}

	ctx.JSON(http.StatusCreated, response.NewLeaseChargePaymentResponse(leaseChargePayment))
}

// swagger:route GET /leases/{id}/balance lease getLeaseBalance
//
// Get a lease balance
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// responses:
// 	200: LeaseBalanceResponse
//  401: ErrorResponse
//  500: ErrorResponse
func (ctrl *leaseController) GetLeaseBalance(ctx *gin.Context) {
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

	var balance int64
	for _, leaseCharge := range *lease.LeaseCharge {
		balance += leaseCharge.Amount
	}

	ctx.JSON(http.StatusOK, response.NewLeaseBalanceResponse(balance))
	// ctx.JSON(http.StatusOK, response.NewLeaseBalanceResponse(lease))
}
