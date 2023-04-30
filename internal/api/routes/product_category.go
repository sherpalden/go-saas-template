package routes

import (
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/controller"
	"github.com/sherpalden/go-saas-template/internal/api/middleware"
)

type ProductCategoryRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	controller     controller.ProductCategoryController
	authMiddleware middleware.AuthMiddleware
}

func NewProductCategoryRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	controller controller.ProductCategoryController,
	authMiddleware middleware.AuthMiddleware,
) ProductCategoryRoutes {
	return ProductCategoryRoutes{
		logger:         logger,
		router:         router,
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (route ProductCategoryRoutes) Setup() {
	route.logger.Zap.Info("Setting up product category routes")

	productCategoryRoutes := route.router.RGroup.Group("/product-categories").Use(route.authMiddleware.Handle)
	productCategoryRoutes.GET("", route.controller.GetCategoryList)
	productCategoryRoutes.POST("", route.controller.Create)
	productCategoryRoutes.DELETE("/:categoryID", route.controller.Create)
}
