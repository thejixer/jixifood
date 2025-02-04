package config

import "os"

type GatewayConfig struct {
	ListenAddr     string
	AuthServiceURI string
	MenuServiceURI string
}

func NewGatewayConfig() *GatewayConfig {

	listenAddr := os.Getenv("GW_PORT")
	AuthServiceURI := os.Getenv("AS_URI")
	MenuServiceURI := os.Getenv("MS_URI")

	return &GatewayConfig{
		ListenAddr:     listenAddr,
		AuthServiceURI: AuthServiceURI,
		MenuServiceURI: MenuServiceURI,
	}
}
