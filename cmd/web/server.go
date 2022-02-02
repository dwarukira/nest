package main

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/solabsafrica/afrikanest/config"
	"github.com/solabsafrica/afrikanest/db"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/queue"
	"github.com/solabsafrica/afrikanest/repo"
	"github.com/solabsafrica/afrikanest/restful"
	"github.com/solabsafrica/afrikanest/restful/middlewares"
	"github.com/solabsafrica/afrikanest/service"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run() error
	Stop() error
}

type restfulServer struct {
	engine *gin.Engine
	config *config.Config
}

func NewServer() (Server, error) {
	config := config.Get()
	engine := gin.New()
	//health-check
	engine.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "Still hanging")
	})
	configs := cors.DefaultConfig()
	configs.AllowAllOrigins = true
	configs.AllowCredentials = true
	configs.AddAllowHeaders("authorization")
	engine.Use(cors.New(configs))

	// engine.Use(cors.Default())
	engine.Use(gin.Logger())

	routerGroup := engine.Group("api")

	setup(routerGroup)

	return &restfulServer{
		engine: engine,
		config: config,
	}, nil
}

func setup(routerGroup *gin.RouterGroup) {
	database := db.NewDatabase()
	redisClient := db.NewRedis()
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("failed to ping redis %v", err)
	}

	userRepo := repo.NewUserRepoWithContext(database)
	propertyRepo := repo.NewPropertyRepoWithContext(database)
	leaseRepo := repo.NewLeaseRepoWithContext(database)
	ticketRepo := repo.NewTicketRepoWithContext(database)
	tenantRepo := repo.NewTenantRepoWithContext(database)
	q := queue.NewJobQueue()

	smsService := service.NewSmsServiceWithContext()
	userService := service.NewUserServiceWithContext(userRepo)
	authService := service.NewAuthServiceWithContext(userRepo)
	propertyService := service.NewPropertyServiceWithContext(propertyRepo)
	leaseService := service.NewLeaseServiceWithContext(leaseRepo)
	ticketService := service.NewTicketServiceWithContext(ticketRepo)
	emailService := service.NewEmailServiceWithContext(q)
	tenantService := service.NewTenantServiceWithContext(tenantRepo, userRepo)

	authChecker := middlewares.NewAuthChecker(authService)

	restful.NewUserController(routerGroup, authChecker, userService)
	restful.NewAuthController(routerGroup, authService)
	restful.NewPropertyController(routerGroup, authChecker, propertyService)
	restful.NewLeaseController(routerGroup, authChecker, leaseService, propertyService, smsService, emailService)
	restful.NewTicketController(routerGroup, authChecker, ticketService)
	restful.NewTenantController(routerGroup, authChecker, tenantService)
}

func (server *restfulServer) Run() error {
	port := server.config.ServerConfig.Port
	return server.engine.Run(":" + port)
}

func (server *restfulServer) Stop() error {
	// clear up any resource
	return nil
}
