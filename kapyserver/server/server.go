package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/k3s-io/k3s/pkg/agent"
	"github.com/k3s-io/k3s/pkg/cli/cmds"
	"github.com/k3s-io/k3s/pkg/clientaccess"
	"github.com/k3s-io/k3s/pkg/server"
	"github.com/kapycluster/corpy/kapyserver/config"
	kcpb "github.com/kapycluster/corpy/types/kubeconfig"
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

	errCh := make(chan error)
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		if err := run(ctx, serverConfig); err != nil {
			errCh <- fmt.Errorf("running k3s server: %w", err)
		}
	}()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	kcpb.RegisterKubeConfigServiceServer(grpcServer, &kubeConfigServer{
		config: serverConfig,
	})
	go func() {
		wg.Add(1)
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("running kubeconfig gRPC server: %w", err)
		}
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			return err
		}
	case sig := <-signals:
		fmt.Printf("recieved signal: %s\n", sig)
		cancel()
		return nil
	}

	<-errCh
	return nil
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
