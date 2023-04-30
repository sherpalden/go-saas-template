package model

import (
	"time"

	"github.com/sherpalden/go-saas-template/external/password"
	"github.com/sherpalden/go-saas-template/internal/entity/admin"
	"gorm.io/gorm"
)

type Admin struct {
	ID        ID             `json:"id"`
	Name      string         `json:"name" validate:"required"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"password" validate:"required"`
	Role      admin.Role     `json:"role" validate:"required,role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m Admin) TableName() string {
	return "admins"
}

func (m *Admin) BeforeCreate(db *gorm.DB) error {
	passwordHash, err := password.GenerateHash(m.Password)
	if err != nil {
		return err
	}

	m.ID = NewID()
	m.Password = passwordHash

	return nil
}
