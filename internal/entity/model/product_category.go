package model

import (
	"time"

	"gorm.io/gorm"
)	

type ProductCategory struct {
	ID          ID        `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name" validate:"required"`
	Lft         int32     `json:"lft"`
	Rgt         int32     `json:"rgt"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m ProductCategory) TableName() string {
	return "product_categories"
}

func (m *ProductCategory) BeforeCreate(db *gorm.DB) error {
	m.ID = NewID()
	return nil
}

type ProductCategoryParams struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

type ProductCategoryWithDepth struct {
	ProductCategory
	Depth int32 `json:"depth"`
}
