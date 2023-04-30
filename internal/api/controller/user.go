package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/external/auth"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/api/service"
	"github.com/sherpalden/go-saas-template/internal/api/validator"
	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
)

type UserController struct {
	userService   service.UserService
	userValidator validator.UserValidator
	logger        infrastructure.Logger
}

func NewUserController(
	userService service.UserService,
	userValidator validator.UserValidator,
	logger infrastructure.Logger,
) UserController {
	return UserController{
		userService:   userService,
		userValidator: userValidator,
		logger:        logger,
	}
}

func (ctrl UserController) Login(ctx *gin.Context) {
	loginCredential := auth.LoginCredential{}
	if err := ctx.ShouldBindJSON(&loginCredential); err != nil {
		err = app_error.New(err, http.StatusInternalServerError)
		httpi.HandleError(ctx, err)
		return
	}

	if validations, passed := validator.GenerateValidation(ctrl.userValidator, loginCredential); !passed {
		err := app_error.New(errors.New("validation error"), http.StatusBadRequest).SetFieldErrors(validations)
		httpi.HandleError(ctx, err)
		return
	}

	tokens, err := ctrl.userService.Login(ctx, loginCredential)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}

	ctx.SetCookie("refresh_token", tokens.RefreshToken, int(auth.REFRESH_TOKEN_VALIDITY_PERIOD), "/", "localhost", false, true)
	httpi.JSON(ctx, http.StatusOK, tokens.AccessToken)
}

func (ctrl UserController) Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	tokens, err := ctrl.userService.Refresh(ctx, refreshToken)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}

	ctx.SetCookie("refresh_token", tokens.RefreshToken, int(auth.REFRESH_TOKEN_VALIDITY_PERIOD), "/", "localhost", false, true)
	httpi.JSON(ctx, http.StatusOK, tokens.AccessToken)
}

func (ctrl UserController) SignUp(ctx *gin.Context) {
	newUser := model.User{}
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		err = app_error.New(err, http.StatusInternalServerError)
		httpi.HandleError(ctx, err)
		return
	}

	if validations, passed := validator.GenerateValidation(ctrl.userValidator, newUser); !passed {
		err := app_error.New(errors.New("validation error"), http.StatusBadRequest).SetFieldErrors(validations)
		httpi.HandleError(ctx, err)
		return
	}

	createdUser, err := ctrl.userService.Create(ctx, newUser)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}

	httpi.JSON(ctx, http.StatusOK, createdUser)
}

func (ctrl UserController) FindAll(ctx *gin.Context) {
	paginationReq := httpi.PaginationRequest{}
	if err := httpi.BuildQueryParams(ctx, paginationReq); err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	users, paginationRes, err := ctrl.userService.FindAll(ctx, paginationReq)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	httpi.JSONWithMetaData(ctx, http.StatusOK, users, paginationRes)
}
