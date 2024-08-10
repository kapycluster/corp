package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/k3s-io/k3s/pkg/agent"
	"github.com/k3s-io/k3s/pkg/cli/cmds"
	"github.com/k3s-io/k3s/pkg/clientaccess"
	"github.com/k3s-io/k3s/pkg/server"
	"github.com/kapycluster/corpy/kapyserver/config"
	"github.com/kapycluster/corpy/kapyserver/util"
	"github.com/kapycluster/corpy/types"
	kcpb "github.com/kapycluster/corpy/types/kubeconfig"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return fmt.Errorf("config failed: %w", err)
	}

	g := errgroup.Group{}

	g.Go(func() error {
		return run(ctx, serverConfig)
	})

	lis, err := net.Listen("tcp", util.MustGetEnv(types.KapyServerGRPCAddress))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	kcpb.RegisterKubeConfigServiceServer(grpcServer, &kubeConfigServer{
		config: serverConfig,
	})
	g.Go(func() error {
		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("running kubeconfig gRPC server: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	select {
	case sig := <-signals:
		fmt.Printf("recieved signal: %s\n", sig)
		cancel()
		return nil
	}
}

// Our own minimal k3s run function
func run(ctx context.Context, serverConfig *config.ServerConfig) error {
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
