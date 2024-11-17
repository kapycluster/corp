package kapyclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"kapycluster.com/corp/types/proto"
)

func (k *KapyClient) GenerateToken(ctx context.Context) (string, error) {
	tokenClient := proto.NewTokenClient(k.client)
	treq := &proto.TokenRequest{}

	var callOpts []grpc.CallOption
	token, err := tokenClient.GenerateToken(ctx, treq, callOpts...)
	if err != nil {
		return "", fmt.Errorf("kapyclient: fetching kubeconfig: %w", err)
	}

	return token.Token, nil
}
