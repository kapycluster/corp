package server

import (
	"github.com/k3s-io/k3s/pkg/cli/cmds"
	"github.com/kapycluster/corpy/kapyserver/config"
)

func Start() error {
	serverConfig := config.NewServerConfig()
	if err := run(serverConfig); err != nil {
		return err
	}

	return nil
}

// Our own minimal k3s run function
func run(config *config.ServerConfig) error {

	if err := cmds.EvacuateCgroup2(); err != nil {
		return err
	}

	// TODO: investigate this and rewrite
	if err := cmds.InitLogging(); err != nil {
		return err
	}

	return nil
}
