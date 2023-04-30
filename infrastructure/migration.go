package infrastructure

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

type Migrations struct {
	logger Logger
	dbConn DBConn
	env    Env
}

func NewMigrations(
	logger Logger,
	dbConn DBConn,
	env Env,
) Migrations {
	return Migrations{
		logger: logger,
		dbConn: dbConn,
		env:    env,
	}
}

func (m Migrations) Migrate() {

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/",
	}

	m.dbConn.DB.Connection(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET ROLE %s", m.env.DBRoleAdmin)).Error
	})

	sqlDB, err := m.dbConn.DB.DB()
	if err != nil {
		m.logger.Zap.Panicf("error getting db connection for migration: %v", err)
	}

	m.logger.Zap.Info("running migration.")

	_, err = migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		m.logger.Zap.Panicf("error executing migration: %v", err)
	}

	m.logger.Zap.Info("migration completed.")
}
