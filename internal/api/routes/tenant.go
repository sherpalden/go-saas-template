package routes

import (
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/controller"
	"github.com/sherpalden/go-saas-template/internal/api/middleware"
)

type TenantRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	controller     controller.TenantController
	authMiddleware middleware.AuthMiddleware
}

func NewTenantRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller controller.TenantController,
	authMiddleware middleware.AuthMiddleware,
) TenantRoutes {
	return TenantRoutes{
		logger:         logger,
		router:         router,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (route TenantRoutes) Setup() {
	route.logger.Zap.Info("Setting up tenant routes")

	tenantRoutes := route.router.RGroup.Group("/tenants").Use(route.authMiddleware.Handle)
	tenantRoutes.POST("", route.controller.Create)

	route.logger.Zap.Info("Setting up tenant users routes")

	tenantUserRoutes := route.router.RGroup.Group("/tenant-users").Use(route.authMiddleware.Handle)
	tenantUserRoutes.POST("", route.controller.AddTenantUser)
	tenantUserRoutes.GET("", route.controller.FindAllTenantUsers)
}
