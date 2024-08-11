package kapyclient

import (
	"context"
	"fmt"

	"github.com/kapycluster/corpy/types/kubeconfig"
	"google.golang.org/grpc"
)

func (k *KapyClient) GetKubeConfig(ctx context.Context) ([]byte, error) {
	kcfgClient := kubeconfig.NewKubeConfigServiceClient(k.client)
	kreq := &kubeconfig.GetKubeConfigRequest{}

	var callOpts []grpc.CallOption
	kcfg, err := kcfgClient.GetKubeConfig(ctx, kreq, callOpts...)
	if err != nil {
		return nil, fmt.Errorf("kapyclient: fetching kubeconfig: %w", err)
	}

	return []byte(kcfg.KubeConfig), nil
}
