package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/external/auth"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/api/repository"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
	"github.com/sherpalden/go-saas-template/internal/entity/user"
)

type TenantService struct {
	db               repository.Database
	tenantRepository repository.TenantRepository
	authClient       auth.AuthClient
	logger           infrastructure.Logger
}

func NewTenantService(
	db repository.Database,
	tenantRepository repository.TenantRepository,
	authClient auth.AuthClient,
	logger infrastructure.Logger,
) TenantService {
	return TenantService{
		db:               db,
		tenantRepository: tenantRepository,
		authClient:       authClient,
		logger:           logger,
	}
}

func (service TenantService) Create(ctx *gin.Context, tenant model.Tenant) (*model.TenantUser, error) {
	tenantUser := &model.TenantUser{}
	if err := service.db.WithAdmin().Atomic(func(db repository.Database) error {
		createdTenant, err := service.tenantRepository.CreateTenant(db, tenant)
		if err != nil {
			return err
		}

		authUser := ctx.MustGet("auth_user")
		userID, err := model.StringToID(authUser.(map[string]interface{})["id"].(string))
		if err != nil {
			return err
		}

		tenantSuperUser := model.TenantUser{
			TenantID:      createdTenant.TenantID,
			UserID:        userID,
			Role:          user.SuperUser,
			AccountStatus: user.Active,
		}
		tenantUser, err = service.tenantRepository.CreateTenantSuperUser(db, tenantSuperUser)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tenantUser, nil
}

func (service TenantService) AddTenantUser(ctx *gin.Context, tenantUser model.TenantUser) (*model.TenantUser, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, err
	}
	return service.tenantRepository.AddTenantUser(db, tenantUser)
}

func (service TenantService) FindAllTenantUsers(ctx *gin.Context, pagination httpi.PaginationRequest) (*[]model.TenantUser, httpi.PaginationResponse, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, httpi.PaginationResponse{}, err
	}
	return service.tenantRepository.FindAllTenantUsers(db, pagination)
}

func (service TenantService) FindOneByID(ctx *gin.Context, ID model.ID) (*model.TenantUser, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	userID, err := model.StringToID(ctx.MustGet("auth_user").(map[string]interface{})["id"].(string))
	if err != nil {
		return nil, err
	}
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, err
	}

	return service.tenantRepository.FindOneByID(db, userID)
}
