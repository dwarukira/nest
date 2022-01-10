package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/restful/middlewares"
	"github.com/solabsafrica/afrikanest/restful/request"
	"github.com/solabsafrica/afrikanest/restful/response"
	"github.com/solabsafrica/afrikanest/service"
)

type ticketController struct {
	ticketService service.TicketServiceWithContext
}

func NewTicketController(group *gin.RouterGroup, authChecker middlewares.AuthChecker, ticketService service.TicketServiceWithContext) {
	controller := &ticketController{
		ticketService: ticketService,
	}

	v1 := group.Group("/v1")
	v1.POST("/tickets", authChecker.Check, controller.CreateTicketHandler)
	// v1.GET("/tickets", authChecker.Check, controller.GetTicketsHandler)
	// v1.GET("/tickets/:id", authChecker.Check, controller.GetTicketHandler)
}

func (ctrl *ticketController) CreateTicketHandler(ctx *gin.Context) {
	var createTicketRequest request.CreateTicketRequest
	if err := ctx.ShouldBindJSON(&createTicketRequest); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	ticket, err := createTicketRequest.ToTicket()
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	if ticket, err = ctrl.ticketService(ctx).Create(ticket); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, response.NewCreateTicketResponse(ticket))
	}
}
