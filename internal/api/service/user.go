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

type UserService struct {
	db               repository.Database
	userRepository   repository.UserRepository
	tenantRepository repository.TenantRepository
	authClient       auth.AuthClient
	logger           infrastructure.Logger
}

func NewUserService(
	db repository.Database,
	userRepository repository.UserRepository,
	tenantRepository repository.TenantRepository,
	authClient auth.AuthClient,
	logger infrastructure.Logger,
) UserService {
	return UserService{
		db:               db,
		userRepository:   userRepository,
		tenantRepository: tenantRepository,
		authClient:       authClient,
		logger:           logger,
	}
}

func (service UserService) Create(ctx *gin.Context, user model.User) (*model.User, error) {
	return service.userRepository.Create(service.db.WithAdmin(), user)
}

func (service UserService) FindAll(ctx *gin.Context, pagination httpi.PaginationRequest) (*[]model.User, httpi.PaginationResponse, error) {
	return service.userRepository.FindAll(service.db.WithAdmin(), pagination)
}

func (service UserService) FindOneByID(ctx *gin.Context, ID model.ID) (*model.User, error) {
	return service.userRepository.FindOneByID(service.db.WithAdmin(), ID)
}

func (service UserService) Login(ctx *gin.Context, loginCredential auth.LoginCredential) (*auth.AuthTokens, error) {
	user, err := service.userRepository.FindOneByEmail(service.db.WithAdmin(), loginCredential.Email)
	if err != nil {
		return nil, err
	}
	if err = password.Verify(user.Password, loginCredential.Password); err != nil {
		return nil, errors.New("incorrect password")
	}

	tenant_id := ""
	role := ""

	if ctx.Query("tenant_id") != "" {
		tenant_id = ctx.Query("tenant_id")
		db, err := service.db.WithRole(tenant_id)
		if err != nil {
			return nil, err
		}
		tenantUser, err := service.tenantRepository.FindOneByID(db, user.ID)
		if err != nil {
			return nil, err
		}
		role = string(tenantUser.Role)
	}

	claims := map[string]interface{}{
		"is_admin":  false,
		"role":      role,
		"tenant_id": tenant_id,
		"id":        user.ID.String(),
	}
	tokens, err := service.authClient.GenerateTokens(claims)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (service UserService) Refresh(ctx *gin.Context, refreshToken string) (*auth.AuthTokens, error) {
	authUser, err := service.authClient.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	userID, err := model.StringToID(authUser["id"].(string))
	if err != nil {
		return nil, err
	}
	user, err := service.userRepository.FindOneByID(service.db, userID)
	if err != nil {
		return nil, err
	}

	role := authUser["role"]
	if authUser["tenant_id"] != "" {
		db, err := service.db.WithRole(authUser["tenant_id"].(string))
		if err != nil {
			return nil, err
		}
		tenantUser, err := service.tenantRepository.FindOneByID(db, user.ID)
		if err != nil {
			return nil, err
		}
		role = string(tenantUser.Role)
	}

	claims := map[string]interface{}{
		"is_admin":  false,
		"role":      role,
		"tenant_id": authUser["tenant_id"],
		"id":        user.ID.String(),
	}
	tokens, err := service.authClient.GenerateTokens(claims)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
