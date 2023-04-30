package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/repository"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
)

type ProductCategoryService struct {
	db                        repository.Database
	productCategoryRepository repository.ProductCategoryRepository
	logger                    infrastructure.Logger
}

func NewProductCategoryService(
	db repository.Database,
	productCategoryRepository repository.ProductCategoryRepository,
	logger infrastructure.Logger,
) ProductCategoryService {
	return ProductCategoryService{
		db:                        db,
		productCategoryRepository: productCategoryRepository,
		logger:                    logger,
	}
}

func (service ProductCategoryService) CreateRootCategory(
	ctx *gin.Context,
	pc model.ProductCategory,
) (*model.ProductCategory, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, err
	}
	pc.TenantID = tenant_id
	return service.productCategoryRepository.CreateRootNode(db, pc)
}

func (service ProductCategoryService) AddNewCategory(
	ctx *gin.Context,
	parentNodeID model.ID,
	pc model.ProductCategory,
) (*model.ProductCategory, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, err
	}
	pc.TenantID = tenant_id
	return service.productCategoryRepository.AddNewNode(db, parentNodeID, pc)
}

func (service ProductCategoryService) DeleteCategory(
	ctx *gin.Context,
	categoryID model.ID,
) error {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return err
	}
	return service.productCategoryRepository.DeleteSubTree(db, categoryID)
}

func (service ProductCategoryService) GetCategoryList(
	ctx *gin.Context,
	params model.ProductCategoryParams,
) ([]model.ProductCategoryWithDepth, error) {
	tenant_id := ctx.MustGet("auth_user").(map[string]interface{})["tenant_id"].(string)
	db, err := service.db.WithRole(tenant_id)
	if err != nil {
		return nil, err
	}
	return service.productCategoryRepository.GetSubTree(db, params)
}
