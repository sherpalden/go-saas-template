package repository

import (
	"fmt"

	"github.com/sherpalden/go-saas-template/infrastructure"
	"gorm.io/gorm"
)

type Database struct {
	conn   infrastructure.DBConn
	env    infrastructure.Env
	logger infrastructure.Logger
}

func NewDatabase(
	conn infrastructure.DBConn,
	env infrastructure.Env,
	logger infrastructure.Logger,
) Database {
	return Database{
		conn:   conn,
		env:    env,
		logger: logger,
	}
}

func (db Database) Atomic(
	fn func(db Database) error,
) (err error) {
	db.logger.Zap.Info("beginning database transaction")
	tx := db.conn.Begin()
	defer func() {
		if p := recover(); p != nil {
			db.logger.Zap.Infof("rolling back database transaction due to: %v", p)
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			db.logger.Zap.Infof("rolling back database transaction due to: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				db.logger.Zap.Errorf("tx err: %v, rb err: %v", err, rbErr)
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			db.logger.Zap.Info("commiting database transaction")
			err = tx.Commit().Error
		}
	}()
	db.conn.DB = tx
	err = fn(db)
	return
}

func (db Database) WithRole(role string) (Database, error) {
	if err := db.conn.DB.Connection(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET ROLE %s", role)).Error
	}); err != nil {
		db.logger.Zap.Errorf("failed to set database user role to %v due to ", role, err)
		return db, err
	}
	return db, nil
}

func (db Database) WithAdmin() Database {
	if err := db.conn.DB.Connection(func(tx *gorm.DB) error {
		return tx.Exec(fmt.Sprintf("SET ROLE %s", db.env.DBRoleAdmin)).Error
	}); err != nil {
		db.logger.Zap.Errorf("failed to set database user role to %v due to ", db.env.DBRoleAdmin, err)
	}
	return db
}
