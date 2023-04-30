package model

import (
	"time"

	"github.com/sherpalden/go-saas-template/internal/entity/user"
	"gorm.io/gorm"
)

type TenantUser struct {
	ID            ID                 `json:"id"`
	TenantID      string             `json:"tenant_id" validate:"required"`
	Tenant        Tenant             `gorm:"References:tenant_id" json:"tenant" validate:"-"`
	UserID        ID                 `json:"user_id" validate:"required"`
	User          User               `json:"user" validate:"-"`
	Role          user.Role          `json:"role" validate:"required"`
	AccountStatus user.AccountStatus `json:"account_status" validate:"required"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
}

func (m TenantUser) TableName() string {
	return "tenants_users"
}

func (m *TenantUser) BeforeCreate(db *gorm.DB) error {
	m.ID = NewID()
	return nil
}
