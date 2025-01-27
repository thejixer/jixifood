package config

import "os"

type GatewayConfig struct {
	ListenAddr string
}

func NewGatewayConfig() *GatewayConfig {

	listenAddr := os.Getenv("GW_PORT")

	return &GatewayConfig{
		ListenAddr: listenAddr,
	}
}
