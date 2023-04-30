package infrastructure

import (
	"os"
)

type Env struct {
	Environment string
	ServerPort  string
	LogOutput   string

	DBRoleAdmin           string
	DBRoleTenant           string
	DBUserGeneral         string
	DBUserGeneralPassword string
	DBHost                string
	DBPort                string
	DBName                string

	SuperAdminName     string
	SuperAdminEmail    string
	SuperAdminPassword string

	AccessTokenSecret  string
	RefreshTokenSecret string
}

func NewEnv() Env {
	env := Env{
		Environment: os.Getenv("ENVIRONMENT"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		LogOutput:   os.Getenv("LOG_OUTPUT"),

		DBRoleAdmin:           os.Getenv("DB_ROLE_ADMIN"),
		DBRoleTenant:           os.Getenv("DB_ROLE_TENANT"),
		DBUserGeneral:         os.Getenv("DB_USER_GENERAL"),
		DBUserGeneralPassword: os.Getenv("DB_USER_GENERAL_PASSWORD"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBName:                os.Getenv("DB_NAME"),

		SuperAdminName:     os.Getenv("SUPER_ADMIN_NAME"),
		SuperAdminEmail:    os.Getenv("SUPER_ADMIN_EMAIL"),
		SuperAdminPassword: os.Getenv("SUPER_ADMIN_PASSWORD"),

		AccessTokenSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
	}

	return env
}
