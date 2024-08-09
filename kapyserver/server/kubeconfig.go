package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kapycluster/corpy/kapyserver/config"
	"github.com/kapycluster/corpy/types/kubeconfig"
)

type kubeConfigServer struct {
	kubeconfig.UnimplementedKubeConfigServiceServer
	config *config.ServerConfig
}

func (k *kubeConfigServer) GetKubeConfig(
	ctx context.Context, gkcr *kubeconfig.GetKubeConfigRequest,
) (*kubeconfig.GetKubeConfigResponse, error) {

	kcfgFile := k.config.ControlConfig.KubeConfigOutput
	kcfg, err := os.ReadFile(kcfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}

	return &kubeconfig.GetKubeConfigResponse{
		KubeConfig: string(kcfg),
	}, nil
}
