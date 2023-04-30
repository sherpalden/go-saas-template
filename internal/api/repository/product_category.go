package repository

import (
	"net/http"

	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
	"github.com/sherpalden/go-saas-template/internal/entity/product_category"
	"gorm.io/gorm"
)

type ProductCategoryRepository struct{}

func NewProductCategoryRepository() ProductCategoryRepository {
	return ProductCategoryRepository{}
}

func (repo ProductCategoryRepository) CreateRootNode(
	db Database,
	productCategory model.ProductCategory,
) (*model.ProductCategory, error) {
	productCategory.Lft = 1
	productCategory.Rgt = 2
	productCategory.Name = product_category.RootName
	if err := db.conn.DB.Create(&productCategory).Error; err != nil {
		return nil, err
	}
	return &productCategory, nil
}

func (repo ProductCategoryRepository) AddNewNode(
	db Database,
	parentNodeID model.ID,
	productCategory model.ProductCategory,
) (*model.ProductCategory, error) {
	var parentCategory model.ProductCategory
	if err := db.conn.DB.Model(&model.ProductCategory{}).Where("id = ?", parentNodeID).First(&parentCategory).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("parent product category not found")
			return nil, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find parent product category")
			return nil, err
		}
	}

	if err := db.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ProductCategory{}).Where("lft > ?", parentCategory.Lft).Update("lft", gorm.Expr("lft + ?", 2)).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.ProductCategory{}).Where("rgt > ?", parentCategory.Lft).Update("rgt", gorm.Expr("rgt + ?", 2)).Error; err != nil {
			return err
		}

		productCategory.Lft = parentCategory.Lft + 1
		productCategory.Rgt = parentCategory.Lft + 2
		if err := tx.Create(&productCategory).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &productCategory, nil
}

func (repo ProductCategoryRepository) DeleteSubTree(
	db Database,
	nodeID model.ID,
) error {
	var category model.ProductCategory
	if err := db.conn.DB.Model(&model.ProductCategory{}).Where("id = ?", nodeID).First(&category).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("product category not found")
			return err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find product category")
			return err
		}
	}

	if err := db.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("lft BETWEEN ? AND ?", category.Lft, category.Rgt).
			Delete(&model.ProductCategory{}).Error; err != nil {
			return err
		}

		width := category.Rgt - category.Lft + 1
		if err := tx.Model(&model.ProductCategory{}).Where("lft > ?", category.Rgt).Update("lft", gorm.Expr("lft - ?", width)).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.ProductCategory{}).Where("rgt > ?", category.Rgt).Update("rgt", gorm.Expr("rgt - ?", width)).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo ProductCategoryRepository) GetSubTree(
	db Database,
	params model.ProductCategoryParams,
) ([]model.ProductCategoryWithDepth, error) {
	if params.Name == "" {
		params.Name = product_category.RootName
	}
	queryBuilder := db.conn.DB.Raw(`
		SELECT
			node.id, 
			node.name, 
			node.tenant_id,
			node.lft,
			node.rgt,
			node.description, 
			node.created_at,
			node.updated_at,
			(COUNT(parent.name)-1) AS depth
		FROM  
			product_categories as node,
			product_categories as parent
		WHERE 
			node.lft BETWEEN parent.lft AND parent.rgt
		GROUP BY node.id
		ORDER BY depth
	`)

	productCategoryList := []model.ProductCategoryWithDepth{}
	err := queryBuilder.
		Find(&productCategoryList).Error
	if err != nil {
		return nil, err
	}
	return productCategoryList, nil
}
