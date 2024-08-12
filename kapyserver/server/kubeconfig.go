package server

import (
	"context"
	"fmt"
	"os"

	"github.com/kapycluster/corpy/kapyserver/config"
	"github.com/kapycluster/corpy/types/proto"
)

type kubeConfigServer struct {
	proto.UnimplementedKubeConfigServiceServer
	config *config.ServerConfig
}

func (k *kubeConfigServer) GetKubeConfig(
	ctx context.Context, gkcr *proto.GetKubeConfigRequest,
) (*proto.GetKubeConfigResponse, error) {

	kcfgFile := k.config.ControlConfig.KubeConfigOutput
	kcfg, err := os.ReadFile(kcfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read proto file: %w", err)
	}

	return &proto.GetKubeConfigResponse{
		KubeConfig: string(kcfg),
	}, nil
}
