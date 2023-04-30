package repository

import (
	"net/http"

	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() UserRepository {
	return UserRepository{}
}

func (repo UserRepository) Create(db Database, user model.User) (*model.User, error) {
	err := db.conn.DB.Create(&user).Error
	return &user, err
}

func (repo UserRepository) FindAll(db Database, pagination httpi.PaginationRequest) (*[]model.User, httpi.PaginationResponse, error) {
	users := []model.User{}
	queryBuilder := db.conn.DB.Model(&model.User{})

	paginationResponse := httpi.PaginationResponse{}
	if err := queryBuilder.Count(&paginationResponse.TotalCount).Error; err != nil {
		return &users, httpi.PaginationResponse{}, err
	}

	limit := pagination.PageSize
	offset := (pagination.Page - 1) * pagination.PageSize
	queryBuilder = queryBuilder.Offset(int(offset)).Order(pagination.Sort)
	if !pagination.All {
		queryBuilder = queryBuilder.Limit(int(limit))
	}

	err := queryBuilder.
		Find(&users).Error
	if err != nil {
		return nil, httpi.PaginationResponse{}, err
	}

	paginationResponse.CurrentCount = int64(len(users))
	if paginationResponse.TotalCount%offset == 0 {
		paginationResponse.TotalPages = paginationResponse.TotalCount / offset
	} else {
		paginationResponse.TotalPages = (paginationResponse.TotalCount / offset) + 1
	}
	paginationResponse.CurrentPage = pagination.Page

	return &users, paginationResponse, nil
}

func (repo UserRepository) FindOneByID(db Database, ID model.ID) (*model.User, error) {
	user := model.User{}
	err := db.conn.DB.Model(&model.User{}).Where("id = ?", ID).First(&user).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("user not found")
			return &user, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find user")
			return &user, err
		}
	}
	return &user, err
}

func (repo UserRepository) FindOneByEmail(db Database, email string) (*model.User, error) {
	user := model.User{}
	err := db.conn.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("user not found")
			return &user, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find user")
			return &user, err
		}
	}
	return &user, err
}
