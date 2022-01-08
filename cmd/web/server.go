package main

import (
	"context"
	"net/http"

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
		logger.Errorf("failde to ping redis %v", err)
	}

	userRepo := repo.NewUserRepoWithContext(database)
	propertyRepo := repo.NewPropertyRepoWithContext(database)
	_ = queue.NewJobQueue()

	userService := service.NewUserServiceWithContext(userRepo)
	authService := service.NewAuthServiceWithContext(userRepo)
	propertyService := service.NewPropertyServiceWithContext(propertyRepo)
	authChecker := middlewares.NewAuthChecker(authService)

	restful.NewUserController(routerGroup, authChecker, userService)
	restful.NewAuthController(routerGroup, authService)
	restful.NewPropertyController(routerGroup, authChecker, propertyService)
}

func (server *restfulServer) Run() error {
	port := server.config.ServerConfig.Port
	return server.engine.Run(":" + port)
}

func (server *restfulServer) Stop() error {
	// clear up any resource
	return nil
}
