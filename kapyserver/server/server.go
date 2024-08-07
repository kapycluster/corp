package server

import (
	"fmt"
	"log"

	"github.com/k3s-io/k3s/pkg/agent"
	"github.com/k3s-io/k3s/pkg/cli/cmds"
	"github.com/k3s-io/k3s/pkg/clientaccess"
	"github.com/k3s-io/k3s/pkg/server"
	"github.com/kapycluster/corpy/kapyserver/config"
	"github.com/rancher/wrangler/v3/pkg/signals"
)

func Start() error {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return fmt.Errorf("config failed: %w", err)
	}

	if err := run(serverConfig); err != nil {
		return err
	}

	return nil
}

// Our own minimal k3s run function
func run(serverConfig *config.ServerConfig) error {
	ctx := signals.SetupSignalContext()

	// TODO: investigate this and rewrite

	cmds.LogConfig.VLevel = 5
	if err := cmds.InitLogging(); err != nil {
		return fmt.Errorf("initializing logging: %w", err)
	}

	// We have to pass an empty cmds.Server{} since it's needed down below
	if err := server.StartServer(ctx, &serverConfig.Config, &cmds.Server{}); err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	token, err := clientaccess.FormatToken(serverConfig.ControlConfig.Runtime.AgentToken, serverConfig.ControlConfig.Runtime.ServerCA)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Println(token)

	agentConfig := cmds.Agent{
		DataDir:                  serverConfig.ControlConfig.DataDir,
		Debug:                    true,
		ServerURL:                fmt.Sprintf("https://%s:%d", serverConfig.LBAddress, serverConfig.SupervisorPort),
		Token:                    token,
		ContainerRuntimeReady:    config.ContainerRuntimeReady,
		ContainerRuntimeEndpoint: "/dev/null",
		DisableLoadBalancer:      true,
		Rootless:                 true,
	}

	go func() {
		<-serverConfig.ControlConfig.Runtime.APIServerReady
		log.Println("apiserver is up")
		<-serverConfig.ControlConfig.Runtime.ETCDReady
		log.Println("etcd is up")
	}()

	return agent.RunStandalone(ctx, agentConfig)
}
