package config

import "os"

type GatewayConfig struct {
	ListenAddr     string
	AuthServiceURI string
}

func NewGatewayConfig() *GatewayConfig {

	listenAddr := os.Getenv("GW_PORT")
	AuthServiceURI := os.Getenv("AS_URI")

	return &GatewayConfig{
		ListenAddr:     listenAddr,
		AuthServiceURI: AuthServiceURI,
	}
}
