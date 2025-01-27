package config

import "os"

type AuthServiceConfig struct {
	ListenAddr  string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	RedisURI    string
}

func NewAuthServiceConfig() *AuthServiceConfig {

	listenAddr := os.Getenv("AS_PORT")

	DB_NAME := os.Getenv("AS_DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	RedisURI := os.Getenv("REDIS_URI")

	return &AuthServiceConfig{
		ListenAddr:  listenAddr,
		DB_NAME:     DB_NAME,
		DB_USER:     dbUser,
		DB_PASSWORD: dbPassword,
		DB_HOST:     dbHost,
		RedisURI:    RedisURI,
	}
}
