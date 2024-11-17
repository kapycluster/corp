package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/k3s-io/k3s/pkg/agent"
	"github.com/k3s-io/k3s/pkg/cli/cmds"
	"github.com/k3s-io/k3s/pkg/clientaccess"
	"github.com/k3s-io/k3s/pkg/server"
	"google.golang.org/grpc"
	"kapycluster.com/corp/kapyserver/config"
	"kapycluster.com/corp/kapyserver/util"
	"kapycluster.com/corp/log"
	"kapycluster.com/corp/types"
	"kapycluster.com/corp/types/proto"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = log.NewContext(ctx, "kapyserver")
	l := log.FromContext(ctx)

	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return fmt.Errorf("failed to configure: %w", err)
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error, 2)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterKubeConfigServer(grpcServer, &kubeConfigServer{
		config: serverConfig,
	})
	proto.RegisterTokenServer(grpcServer, &tokenServer{
		config: serverConfig,
	})

	runWg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := run(ctx, serverConfig, runWg)
		if err != nil {
			l.Error("k3s error", "error", err)
			errCh <- err
		}
	}()

	l.Info("waiting for k3s to come up...")
	runWg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", util.GetEnv(types.KapyServerGRPCAddress))
		if err != nil {
			l.Error("failed to listen", "error", err)
			errCh <- err
		}
		defer lis.Close()
		l.Info("starting gRPC server", "address", util.GetEnv(types.KapyServerGRPCAddress))
		err = grpcServer.Serve(lis)
		if err != nil {
			l.Error("grpc error", "error", err)
			errCh <- err
		}
	}()

	// Create a channel to receive signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for an error, signal, or context cancellation
	select {
	case err := <-errCh:
		return err
	case <-sigChan:
		fmt.Println("recieved signal, shutting down...")
		cancel()
		grpcServer.GracefulStop()
	case <-ctx.Done():
		fmt.Println("context canceled, shutting down...")
		grpcServer.GracefulStop()
	}

	wg.Wait()
	return nil
}

// Our own minimal k3s run function
func run(ctx context.Context, serverConfig *config.ServerConfig, wg *sync.WaitGroup) error {
	// Waiting for APIServerReady and ETCDReady
	wg.Add(2)

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
		log.FromContext(ctx).Info("apiserver is up")
		wg.Done()

		<-serverConfig.ControlConfig.Runtime.ETCDReady
		log.FromContext(ctx).Info("etcd is up")
		wg.Done()
	}()

	if err := agent.RunStandalone(ctx, agentConfig); err != nil {
		return err
	}

	return nil
}
