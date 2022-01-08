package restful

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/solabsafrica/afrikanest/exceptions"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/restful/middlewares"
	"github.com/solabsafrica/afrikanest/restful/request"
	"github.com/solabsafrica/afrikanest/restful/response"
	"github.com/solabsafrica/afrikanest/service"
)

type propertyController struct {
	propertyService service.PropertyServiceWithContext
}

func NewPropertyController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, propertyService service.PropertyServiceWithContext) {
	controller := &propertyController{
		propertyService: propertyService,
	}

	v1 := group.Group("/v1")
	v1.POST("/properties", authChecker.Check, controller.CreatePropertyHandler)
	v1.GET("/properties", authChecker.Check, controller.GetPropertiesHandler)
	v1.GET("/units/:id", authChecker.Check, controller.GetUnitHandler)
	v1.POST("/units", authChecker.Check, controller.CreateUnitHandler)
}

func (ctrl *propertyController) CreatePropertyHandler(ctx *gin.Context) {
	var createPropertyRequest request.CreatePropertyRequest
	if err := ctx.ShouldBindJSON(&createPropertyRequest); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	property, err := createPropertyRequest.ToProperty()
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	id, err := GetIndentityFromContext(ctx)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	property.OwnerID = id

	if property, err = ctrl.propertyService(ctx).Create(property); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, response.NewCreatePropertyResponse(property))
	}
}

func (ctrl *propertyController) GetPropertiesHandler(ctx *gin.Context) {
	page := request.NewPageRequest(ctx.Query("page"), ctx.Query("pageSize"))
	i, err := GetIndentityFromContext(ctx)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, err)
	}
	properties, totalCount, err := ctrl.propertyService(ctx).ListUserProperties(page.Page, page.PageSize, i)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, response.NewGetPropertiesResponse(properties, response.NewPagination(page.Page, page.PageSize, totalCount)))
}

func (ctrl *propertyController) GetUnitHandler(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, exceptions.UUIDParseFailed.SetMessage(ctx.Param("id")))
		return
	}
	unit, err := ctrl.propertyService(ctx).GetUnitById(id)
	if err != nil {
		if errors.Is(err, exceptions.UnitNotExists) {
			ctx.JSON(http.StatusNotFound, err)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	property, _ := ctrl.propertyService(ctx).GetPropertyById(unit.PropertyID)
	if err != nil {
		// Units with no valid property can't be fetched
		if errors.Is(err, exceptions.UnitNotExists) {
			ctx.JSON(http.StatusNotFound, err)
		} else {

			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	id, err = GetIndentityFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if property.OwnerID != id {
		ctx.JSON(http.StatusUnauthorized, exceptions.AuthFailed.SetMessage("not authorized"))
		return
	}

	ctx.JSON(http.StatusOK, unit)

}

func (ctrl *propertyController) CreateUnitHandler(ctx *gin.Context) {
	var createUnitRequest request.CreateUnitRequest
	if err := ctx.ShouldBindJSON(&createUnitRequest); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	unit, err := createUnitRequest.ToUnit()
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	property, err := ctrl.propertyService(ctx).GetPropertyById(unit.PropertyID)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	id, err := GetIndentityFromContext(ctx)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if property.ID != unit.PropertyID || property.OwnerID != id {
		err = exceptions.UnitCreateFailed.Wrap(err).SetMessage("failed to create unit for the provided property check access")
		logger.Error(err)

		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	if unit, err = ctrl.propertyService(ctx).CreateUnit(unit); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, response.NewCreateUnitResponse(unit))
	}

}
