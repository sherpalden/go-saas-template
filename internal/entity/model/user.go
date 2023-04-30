package model

import (
	"time"

	"github.com/sherpalden/go-saas-template/external/password"
	"gorm.io/gorm"
)

type User struct {
	ID        ID             `json:"id"`
	Name      string         `json:"name" validate:"required"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"password" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m User) TableName() string {
	return "users"
}

func (m *User) BeforeCreate(db *gorm.DB) error {
	passwordHash, err := password.GenerateHash(m.Password)
	if err != nil {
		return err
	}

	m.ID = NewID()
	m.Password = passwordHash

	return nil
}
