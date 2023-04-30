package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/api/service"
	"github.com/sherpalden/go-saas-template/internal/api/validator"
	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
)

type TenantController struct {
	tenantService   service.TenantService
	tenantValidator validator.TenantValidator
	logger          infrastructure.Logger
}

func NewTenantController(
	tenantService service.TenantService,
	tenantValidator validator.TenantValidator,
	logger infrastructure.Logger,
) TenantController {
	return TenantController{
		tenantService:   tenantService,
		tenantValidator: tenantValidator,
		logger:          logger,
	}
}

func (ctrl TenantController) Create(ctx *gin.Context) {
	tenant := model.Tenant{}
	if err := ctx.ShouldBindJSON(&tenant); err != nil {
		err = app_error.New(err, http.StatusInternalServerError)
		httpi.HandleError(ctx, err)
		return
	}

	if validations, passed := validator.GenerateValidation(ctrl.tenantValidator, tenant); !passed {
		err := app_error.New(errors.New("validation error"), http.StatusBadRequest).SetFieldErrors(validations)
		httpi.HandleError(ctx, err)
		return
	}

	tenantUser, err := ctrl.tenantService.Create(ctx, tenant)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}

	httpi.JSON(ctx, http.StatusOK, tenantUser)
}

func (ctrl TenantController) AddTenantUser(ctx *gin.Context) {
	tenantUser := model.TenantUser{}
	if err := ctx.ShouldBindJSON(&tenantUser); err != nil {
		err = app_error.New(err, http.StatusInternalServerError)
		httpi.HandleError(ctx, err)
		return
	}

	if validations, passed := validator.GenerateValidation(ctrl.tenantValidator, tenantUser); !passed {
		err := app_error.New(errors.New("validation error"), http.StatusBadRequest).SetFieldErrors(validations)
		httpi.HandleError(ctx, err)
		return
	}

	addedTenantUser, err := ctrl.tenantService.AddTenantUser(ctx, tenantUser)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}

	httpi.JSON(ctx, http.StatusOK, addedTenantUser)

}

func (ctrl TenantController) FindAllTenantUsers(ctx *gin.Context) {
	paginationReq := httpi.PaginationRequest{}
	if err := httpi.BuildQueryParams(ctx, &paginationReq); err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	tenantUsers, paginationRes, err := ctrl.tenantService.FindAllTenantUsers(ctx, paginationReq)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	httpi.JSONWithMetaData(ctx, http.StatusOK, tenantUsers, paginationRes)
}

func (ctrl TenantController) FindOneByID(ctx *gin.Context) {
	userID, err := model.StringToID(ctx.MustGet("auth_user").(map[string]interface{})["id"].(string))
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	tenantUser, err := ctrl.tenantService.FindOneByID(ctx, userID)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	httpi.JSON(ctx, http.StatusOK, tenantUser)
}
