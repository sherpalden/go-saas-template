package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/external/auth"
	"github.com/sherpalden/go-saas-template/external/password"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/api/repository"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
)

type AdminService struct {
	db              repository.Database
	adminRepository repository.AdminRepository
	authClient      auth.AuthClient
	logger          infrastructure.Logger
	env             infrastructure.Env
}

func NewAdminService(
	db repository.Database,
	adminRepository repository.AdminRepository,
	authClient auth.AuthClient,
	logger infrastructure.Logger,
	env infrastructure.Env,
) AdminService {
	return AdminService{
		db:              db,
		adminRepository: adminRepository,
		authClient:      authClient,
		logger:          logger,
		env:             env,
	}
}

func (service AdminService) Create(ctx *gin.Context, admin model.Admin) (*model.Admin, error) {
	db := service.db.WithAdmin()
	return service.adminRepository.Create(db, admin)
}

func (service AdminService) FindAll(ctx *gin.Context, pagination httpi.PaginationRequest) (*[]model.Admin, httpi.PaginationResponse, error) {
	db := service.db.WithAdmin()
	return service.adminRepository.FindAll(db, pagination)
}

func (service AdminService) FindOneByID(ctx *gin.Context, ID model.ID) (*model.Admin, error) {
	db := service.db.WithAdmin()
	return service.adminRepository.FindOneByID(db, ID)
}

func (service AdminService) Login(ctx *gin.Context, loginCredential auth.LoginCredential) (*auth.AuthTokens, error) {
	db := service.db.WithAdmin()
	admin, err := service.adminRepository.FindOneByEmail(db, loginCredential.Email)
	if err != nil {
		return nil, err
	}
	if err = password.Verify(admin.Password, loginCredential.Password); err != nil {
		return nil, errors.New("incorrect password")
	}

	claims := map[string]interface{}{
		"is_admin": true,
		"role":     admin.Role,
		"id":       admin.ID.String(),
	}
	tokens, err := service.authClient.GenerateTokens(claims)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (service AdminService) Refresh(ctx *gin.Context, refreshToken string) (*auth.AuthTokens, error) {
	authUser, err := service.authClient.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	adminID, err := model.StringToID(authUser["id"].(string))
	if err != nil {
		return nil, err
	}
	db := service.db.WithAdmin()
	admin, err := service.adminRepository.FindOneByID(db, adminID)
	if err != nil {
		return nil, err
	}

	claims := map[string]interface{}{
		"is_admin": true,
		"role":     admin.Role,
		"id":       admin.ID.String(),
	}
	tokens, err := service.authClient.GenerateTokens(claims)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
