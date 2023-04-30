package seeds

import (
	"github.com/sherpalden/go-saas-template/infrastructure"
	"github.com/sherpalden/go-saas-template/internal/api/repository"
	"github.com/sherpalden/go-saas-template/internal/entity/admin"
	"github.com/sherpalden/go-saas-template/internal/entity/model"
)

type SuperAdminSeed struct {
	db              repository.Database
	adminRepository repository.AdminRepository
	env             infrastructure.Env
	logger          infrastructure.Logger
}

// NewAdminSeed creates admin seed
func NewAdminSeed(
	db repository.Database,
	adminRepository repository.AdminRepository,
	env infrastructure.Env,
	logger infrastructure.Logger,
) SuperAdminSeed {
	return SuperAdminSeed{
		db:              db,
		adminRepository: adminRepository,
		env:             env,
		logger:          logger,
	}
}

func (seeder SuperAdminSeed) Run() {
	superAdmin := model.Admin{
		Name:     seeder.env.SuperAdminName,
		Email:    seeder.env.SuperAdminEmail,
		Password: seeder.env.SuperAdminPassword,
		Role:     admin.SuperAdmin,
	}
	superAdminExists, err := seeder.adminRepository.FindOneByEmail(seeder.db.WithAdmin(), superAdmin.Email)
	if superAdminExists != nil && err == nil {
		seeder.logger.Zap.Info("super admin seed already")
	} else {
		_, err := seeder.adminRepository.Create(seeder.db.WithAdmin(), superAdmin)
		if err != nil {
			seeder.logger.Zap.Warn("unable to seed super admin")
		} else {
			seeder.logger.Zap.Info("super admin seed successful")
		}
	}
}
