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

type ProductCategoryController struct {
	pcService   service.ProductCategoryService
	pcValidator validator.ProductCategoryValidator
	logger      infrastructure.Logger
}

func NewProductCategoryController(
	pcService service.ProductCategoryService,
	pcValidator validator.ProductCategoryValidator,
	logger infrastructure.Logger,
) ProductCategoryController {
	return ProductCategoryController{
		pcService:   pcService,
		pcValidator: pcValidator,
		logger:      logger,
	}
}

func (ctrl ProductCategoryController) Create(ctx *gin.Context) {
	category := model.ProductCategory{}
	if err := ctx.ShouldBindJSON(&category); err != nil {
		err = app_error.New(err, http.StatusInternalServerError)
		httpi.HandleError(ctx, err)
		return
	}

	if validations, passed := validator.GenerateValidation(ctrl.pcValidator, category); !passed {
		err := app_error.New(errors.New("validation error"), http.StatusBadRequest).SetFieldErrors(validations)
		httpi.HandleError(ctx, err)
		return
	}

	parentIDStr := ctx.Query("parent_id")
	if parentIDStr == "" {
		rootCategory, err := ctrl.pcService.CreateRootCategory(ctx, category)
		if err != nil {
			httpi.HandleError(ctx, err)
			return
		}
		httpi.JSON(ctx, http.StatusOK, rootCategory)
	} else {
		parentID, err := model.StringToID(parentIDStr)
		if err != nil {
			httpi.HandleError(ctx, err)
			return
		}
		newCategory, err := ctrl.pcService.AddNewCategory(ctx, parentID, category)
		if err != nil {
			httpi.HandleError(ctx, err)
			return
		}
		httpi.JSON(ctx, http.StatusOK, newCategory)
	}
}

func (ctrl ProductCategoryController) DeleteCategory(ctx *gin.Context) {
	categoryIDStr := ctx.Param("categoryID")
	categoryID, err := model.StringToID(categoryIDStr)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	if err := ctrl.pcService.DeleteCategory(ctx, categoryID); err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	httpi.JSON(ctx, http.StatusOK, "category delete successful")
}

func (ctrl ProductCategoryController) GetCategoryList(ctx *gin.Context) {
	params := model.ProductCategoryParams{}
	params.Name = ctx.Query("name")

	categoryList, err := ctrl.pcService.GetCategoryList(ctx, params)
	if err != nil {
		httpi.HandleError(ctx, err)
		return
	}
	httpi.JSON(ctx, http.StatusOK, categoryList)
}
