package config

import "os"

type AuthServiceConfig struct {
	ListenAddr string
}

func NewAuthServiceConfig() *AuthServiceConfig {

	listenAddr := os.Getenv("AS_PORT")

	return &AuthServiceConfig{
		ListenAddr: listenAddr,
	}
}
