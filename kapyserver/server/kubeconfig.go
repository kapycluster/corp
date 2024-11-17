package server

import (
	"context"
	"fmt"
	"os"

	"kapycluster.com/corp/kapyserver/config"
	"kapycluster.com/corp/types/proto"
)

type kubeConfigServer struct {
	proto.UnimplementedKubeConfigServer
	config *config.ServerConfig
}

func (k *kubeConfigServer) GetKubeConfig(
	ctx context.Context, kr *proto.KubeConfigRequest,
) (*proto.KubeConfigData, error) {

	kcfgFile := k.config.ControlConfig.KubeConfigOutput
	kcfg, err := os.ReadFile(kcfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read proto file: %w", err)
	}

	return &proto.KubeConfigData{
		KubeConfig: string(kcfg),
	}, nil
}
