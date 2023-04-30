package repository

import (
	"net/http"

	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
	"gorm.io/gorm"
)

type AdminRepository struct{}

func NewAdminRepository() AdminRepository {
	return AdminRepository{}
}

func (repo AdminRepository) Create(db Database, admin model.Admin) (*model.Admin, error) {
	err := db.conn.DB.Create(&admin).Error
	return &admin, err
}

func (repo AdminRepository) FindAll(db Database, paginationReq httpi.PaginationRequest) (*[]model.Admin, httpi.PaginationResponse, error) {
	admins := []model.Admin{}
	queryBuilder := db.conn.DB.Model(&model.Admin{})

	paginationResponse := httpi.PaginationResponse{}
	if err := queryBuilder.Count(&paginationResponse.TotalCount).Error; err != nil {
		return &admins, httpi.PaginationResponse{}, err
	}

	limit := paginationReq.PageSize
	offset := (paginationReq.Page - 1) * paginationReq.PageSize
	queryBuilder = queryBuilder.Offset(int(offset)).Order(paginationReq.Sort)
	if !paginationReq.All {
		queryBuilder = queryBuilder.Limit(int(limit))
	}

	err := queryBuilder.
		Find(&admins).Error
	if err != nil {
		return nil, httpi.PaginationResponse{}, err
	}

	paginationResponse.CurrentCount = int64(len(admins))
	if limit != 0 || paginationResponse.TotalCount != 0 {
		if paginationResponse.TotalCount%limit == 0 {
			paginationResponse.TotalPages = paginationResponse.TotalCount / limit
		} else {
			paginationResponse.TotalPages = (paginationResponse.TotalCount / limit) + 1
		}
	}
	paginationResponse.CurrentPage = paginationReq.Page

	return &admins, paginationResponse, nil
}

func (repo AdminRepository) FindOneByID(db Database, ID model.ID) (*model.Admin, error) {
	admin := model.Admin{}
	err := db.conn.DB.Model(&model.Admin{}).Where("id = ?", ID).First(&admin).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("admin not found")
			return &admin, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find admin")
			return &admin, err
		}
	}
	return &admin, nil
}

func (repo AdminRepository) FindOneByEmail(db Database, email string) (*model.Admin, error) {
	admin := model.Admin{}
	err := db.conn.DB.Model(&model.Admin{}).Where("email = ?", email).First(&admin).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusNotFound).SetMessage("admin not found")
			return &admin, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find admin")
			return &admin, err
		}
	}
	return &admin, err
}
