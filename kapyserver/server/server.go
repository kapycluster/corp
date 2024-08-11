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
	"github.com/kapycluster/corpy/kapyserver/util"
	"github.com/kapycluster/corpy/types"
	"github.com/kapycluster/corpy/types/kubeconfig"
	"google.golang.org/grpc"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return fmt.Errorf("failed to configure: %w", err)
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error, 2)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	kubeconfig.RegisterKubeConfigServiceServer(grpcServer, &kubeConfigServer{
		config: serverConfig,
	})

	runWg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := run(ctx, serverConfig, runWg)
		if err != nil {
			log.Printf("k3s error: %v", err)
			errCh <- err
		}
	}()

	log.Println("waiting for k3s to come up...")
	runWg.Wait()
	wg.Add(1)
	go func() {
		lis, err := net.Listen("tcp", util.GetEnv(types.KapyServerGRPCAddress))
		if err != nil {
			log.Printf("failed to listen: %s", err)
			errCh <- err
		}
		defer lis.Close()
		log.Printf("starting gRPC server on %s", util.GetEnv(types.KapyServerGRPCAddress))
		defer wg.Done()
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("grpc error: %v", err)
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
		wg.Done()

		<-serverConfig.ControlConfig.Runtime.ETCDReady
		log.Println("etcd is up")
		wg.Done()
	}()

	if err := agent.RunStandalone(ctx, agentConfig); err != nil {
		return err
	}

	return nil
}
