package routes

import (
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/controller"
	"github.com/sherpalden/go-saas-template/internal/api/middleware"
)

type UserRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	controller     controller.UserController
	authMiddleware middleware.AuthMiddleware
}

func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller controller.UserController,
	authMiddleware middleware.AuthMiddleware,
) UserRoutes {
	return UserRoutes{
		logger:         logger,
		router:         router,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (route UserRoutes) Setup() {
	route.logger.Zap.Info("Setting up user routes")

	route.router.RGroup.POST("/user-sign-up", route.controller.SignUp)
	route.router.RGroup.POST("/user-login", route.controller.Login)
	route.router.RGroup.POST("/user-refresh", route.controller.Refresh)

	userRoutes := route.router.RGroup.Group("/users").Use(route.authMiddleware.AdminOnly)
	userRoutes.GET("", route.authMiddleware.AdminOnly, route.controller.FindAll)
}
