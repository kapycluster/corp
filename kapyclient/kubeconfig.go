package kapyclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"kapycluster.com/corp/types/proto"
)

func (k *KapyClient) GetKubeConfig(ctx context.Context) ([]byte, error) {
	kcfgClient := proto.NewKubeConfigClient(k.client)
	kreq := &proto.KubeConfigRequest{}

	var callOpts []grpc.CallOption
	kcfg, err := kcfgClient.GetKubeConfig(ctx, kreq, callOpts...)
	if err != nil {
		return nil, fmt.Errorf("kapyclient: fetching kubeconfig: %w", err)
	}

	return []byte(kcfg.KubeConfig), nil
}
