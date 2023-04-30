package bootstrap

import (
	"github.com/sherpalden/go-saas-template/external/auth"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/controller"
	"github.com/sherpalden/go-saas-template/internal/api/middleware"
	"github.com/sherpalden/go-saas-template/internal/api/repository"
	"github.com/sherpalden/go-saas-template/internal/api/routes"
	"github.com/sherpalden/go-saas-template/internal/api/seeds"
	"github.com/sherpalden/go-saas-template/internal/api/service"
	"github.com/sherpalden/go-saas-template/internal/api/validator"
)

func Run() {
	logger := infrastructure.NewLogger()
	env := infrastructure.NewEnv()

	dbConn := infrastructure.NewDBConn(env, logger)
	db := repository.NewDatabase(dbConn, env, logger)

	authClient := auth.NewAuthClient(env.AccessTokenSecret, env.RefreshTokenSecret)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	adminRepository := repository.NewAdminRepository()
	userRepository := repository.NewUserRepository()
	tenantRepository := repository.NewTenantRepository()
	pcRepository := repository.NewProductCategoryRepository()

	adminService := service.NewAdminService(db, adminRepository, authClient, logger, env)
	userService := service.NewUserService(db, userRepository, tenantRepository, authClient, logger)
	tenantService := service.NewTenantService(db, tenantRepository, authClient, logger)
	pcService := service.NewProductCategoryService(db, pcRepository, logger)

	adminValidator := validator.NewAdminValidator()
	userValidator := validator.NewUserValidator()
	tenantValidator := validator.NewTenantValidator()
	pcValidator := validator.NewProductCategoryValidator()

	adminController := controller.NewAdminController(adminService, adminValidator, logger)
	userController := controller.NewUserController(userService, userValidator, logger)
	tenantController := controller.NewTenantController(tenantService, tenantValidator, logger)
	pcController := controller.NewProductCategoryController(pcService, pcValidator, logger)

	router := infrastructure.NewRouter(env)
	adminRoutes := routes.NewAdminRoutes(logger, router, adminController, authMiddleware)
	userRoutes := routes.NewUserRoutes(logger, router, userController, authMiddleware)
	tenantRoutes := routes.NewTenantRoutes(logger, router, tenantController, authMiddleware)
	pcRoutes := routes.NewProductCategoryRoutes(logger, router, pcController, authMiddleware)

	routes := routes.NewRoutes(
		adminRoutes,
		userRoutes,
		tenantRoutes,
		pcRoutes,
	)
	routes.Setup()

	migration := infrastructure.NewMigrations(logger, dbConn, env)
	migration.Migrate()

	superAdminSeed := seeds.NewAdminSeed(db, adminRepository, env, logger)
	seeds := seeds.NewSeeds(superAdminSeed)
	seeds.Run()

	router.Run()
}
