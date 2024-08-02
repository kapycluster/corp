package config

import (
	"github.com/k3s-io/k3s/pkg/server"
	"github.com/kapycluster/corpy/kapyserver/util"
)

type ServerConfig struct {
	server.Config
	LBAddress string
}

// NewServerConfig sets defaults and creates a new ServerConfig
func NewServerConfig() *ServerConfig {
	config := &ServerConfig{}
	config.DisableAgent = true
	config.ControlConfig.Token = util.MustGetEnv("KAPYSERVER_TOKEN")
	config.ControlConfig.JoinURL = util.MustGetEnv("KAPYSERVER_JOIN_URL")
	config.ControlConfig.DataDir = util.MustGetEnv("KAPYSERVER_DATA_DIR")
	config.ControlConfig.KubeConfigOutput = util.MustGetEnv("KAPYSERVER_KUBECONFIG_PATH")
	config.ControlConfig.AdvertiseIP = util.MustGetEnv("KAPYSERVER_ADVERTISE_IP")
	config.LBAddress = util.MustGetEnv("KAPYSERVER_LB_ADDRESS")

	config.ControlConfig.ServerNodeName = "kapy-server"
	config.ControlConfig.SANs = append(
		config.ControlConfig.SANs,
		"127.0.0.1",
		"::1",
		"localhost",
		config.LBAddress,
		config.ControlConfig.AdvertiseIP,
	)

	return config
}
