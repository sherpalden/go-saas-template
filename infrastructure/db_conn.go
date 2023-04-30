package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConn struct {
	*gorm.DB
}

func NewDBConn(env Env, zapLogger Logger) DBConn {
	dbUser := env.DBUserGeneral
	dbPassword := env.DBUserGeneralPassword
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbName := env.DBName

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			},
		),
		&gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		zapLogger.Zap.Error("Error connecting to database: ", err)
		zapLogger.Zap.Panic(err)
	}

	dbConn, err := db.DB()
	if err != nil {
		zapLogger.Zap.Error("Error connecting to database: ", err)
		zapLogger.Zap.Panic(err)
	}

	dbConn.SetConnMaxIdleTime(0)
	dbConn.SetConnMaxLifetime(0)
	dbConn.SetMaxIdleConns(9)
	dbConn.SetMaxOpenConns(18)

	zapLogger.Zap.Info("Database connection established.")

	return DBConn{DB: db}
}
