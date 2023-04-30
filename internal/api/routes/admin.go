package routes

import (
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/controller"
	"github.com/sherpalden/go-saas-template/internal/api/middleware"
)

type AdminRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	controller     controller.AdminController
	authMiddleware middleware.AuthMiddleware
}

func NewAdminRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller controller.AdminController,
	authMiddleware middleware.AuthMiddleware,
) AdminRoutes {
	return AdminRoutes{
		logger:         logger,
		router:         router,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (route AdminRoutes) Setup() {
	route.logger.Zap.Info("Setting up admin routes")

	route.router.RGroup.POST("/admin-login", route.controller.Login)
	route.router.RGroup.POST("/admin-refresh", route.controller.Refresh)

	adminRoutes := route.router.RGroup.Group("/admins").Use(route.authMiddleware.Handle).Use(route.authMiddleware.AdminOnly)
	adminRoutes.GET("", route.controller.FindAll)
}
