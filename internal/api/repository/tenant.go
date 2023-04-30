package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sherpalden/go-saas-template/internal/api/httpi"
	"github.com/sherpalden/go-saas-template/internal/app_error"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
	"github.com/sherpalden/go-saas-template/internal/entity/user"
	"gorm.io/gorm"
)

type TenantRepository struct {
}

func NewTenantRepository() TenantRepository {
	return TenantRepository{}
}

func (repo TenantRepository) CreateTenant(db Database, tenant model.Tenant) (*model.Tenant, error) {

	// query1 := fmt.Sprintf(`
	// 	CREATE ROLE %s WITH
	// 		NOLOGIN
	// 		NOSUPERUSER
	// 		NOINHERIT
	// 		NOCREATEDB
	// 		NOCREATEROLE
	// 		NOREPLICATION
	// 		CONNECTION LIMIT -1;
	// 	GRANT USAGE ON SCHEMA public TO %s;
	// 	GRANT SELECT, UPDATE, INSERT, DELETE ON ALL TABLES IN SCHEMA public TO %s;`, db.env.DBRoleTenant, db.env.DBRoleTenant, db.env.DBRoleTenant)

	query2 := fmt.Sprintf("CREATE ROLE %s WITH NOLOGIN NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION CONNECTION LIMIT -1;", tenant.TenantID)
	query3 := fmt.Sprintf("GRANT %s TO %s;", db.env.DBRoleTenant, tenant.TenantID)
	query4 := fmt.Sprintf("GRANT %s TO %s;", tenant.TenantID, db.env.DBUserGeneral)

	// if err := db.conn.DB.Exec(query1).Error; err != nil {
	// 	return nil, err
	// }
	if err := db.conn.DB.Exec(query2).Error; err != nil {
		return nil, err
	}
	if err := db.conn.DB.Exec(query3).Error; err != nil {
		return nil, err
	}
	if err := db.conn.DB.Exec(query4).Error; err != nil {
		return nil, err
	}

	if err := db.conn.DB.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (repo TenantRepository) CreateTenantSuperUser(db Database, tenantUser model.TenantUser) (*model.TenantUser, error) {
	if tenantUser.Role != user.SuperUser {
		return nil, app_error.New(errors.New("user is not a super user"), http.StatusBadRequest)
	}
	err := db.conn.DB.Omit("Tenant").Omit("User").Create(&tenantUser).Error
	return &tenantUser, err
}

func (repo TenantRepository) AddTenantUser(db Database, tenantUser model.TenantUser) (*model.TenantUser, error) {
	err := db.conn.DB.Omit("Tenant").Omit("User").Create(&tenantUser).Error
	return &tenantUser, err
}

func (repo TenantRepository) FindAllTenantUsers(db Database, pagination httpi.PaginationRequest) (*[]model.TenantUser, httpi.PaginationResponse, error) {
	tenantUsers := []model.TenantUser{}
	queryBuilder := db.conn.DB.Model(&model.TenantUser{})

	paginationResponse := httpi.PaginationResponse{}
	if err := queryBuilder.Count(&paginationResponse.TotalCount).Error; err != nil {
		return &tenantUsers, httpi.PaginationResponse{}, err
	}

	limit := pagination.PageSize
	offset := (pagination.Page - 1) * pagination.PageSize
	queryBuilder = queryBuilder.Offset(int(offset)).Order(pagination.Sort)
	if !pagination.All {
		queryBuilder = queryBuilder.Limit(int(limit))
	}

	err := queryBuilder.
		Preload("Tenant").
		Preload("User").
		Find(&tenantUsers).Error
	if err != nil {
		return nil, httpi.PaginationResponse{}, err
	}

	if limit != 0 || paginationResponse.TotalCount != 0 {
		if paginationResponse.TotalCount%limit == 0 {
			paginationResponse.TotalPages = paginationResponse.TotalCount / limit
		} else {
			paginationResponse.TotalPages = (paginationResponse.TotalCount / limit) + 1
		}
	}
	paginationResponse.CurrentPage = pagination.Page

	return &tenantUsers, paginationResponse, nil
}

func (repo TenantRepository) FindOneByID(db Database, userID model.ID) (*model.TenantUser, error) {
	tenantUser := model.TenantUser{}
	err := db.conn.DB.Model(&model.TenantUser{}).
		Where("user_id = ?", userID).
		Preload("Tenant").
		Preload("User").
		First(&tenantUser).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("user not found")
			return &tenantUser, err
		default:
			err := app_error.New(err, http.StatusInternalServerError).SetMessage("failed to find user")
			return &tenantUser, err
		}
	}
	return &tenantUser, nil
}
