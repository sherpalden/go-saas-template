package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID          ID             `json:"id"`
	TenantID    string         `json:"tenant_id" validate:"required"`
	CompanyName string         `json:"company_name" validate:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m Tenant) TableName() string {
	return "tenants"
}

func (m *Tenant) BeforeCreate(db *gorm.DB) error {
	m.ID = NewID()
	return nil
}
