package config

import "os"

type MenuServiceConfig struct {
	ListenAddr     string
	DB_NAME        string
	DB_USER        string
	DB_PASSWORD    string
	DB_HOST        string
	AuthServiceURI string
}

func NewMenuServiceConfig() *MenuServiceConfig {

	listenAddr := os.Getenv("MS_PORT")

	DB_NAME := os.Getenv("MS_DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	AuthServiceURI := os.Getenv("AS_URI")

	return &MenuServiceConfig{
		ListenAddr:     listenAddr,
		DB_NAME:        DB_NAME,
		DB_USER:        dbUser,
		DB_PASSWORD:    dbPassword,
		DB_HOST:        dbHost,
		AuthServiceURI: AuthServiceURI,
	}
}
