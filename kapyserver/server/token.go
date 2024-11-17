package server

import (
	"context"
	"fmt"

	"github.com/k3s-io/k3s/pkg/clientaccess"
	"kapycluster.com/corp/kapyserver/config"
	"kapycluster.com/corp/types/proto"
)

type tokenServer struct {
	proto.UnimplementedTokenServer
	config *config.ServerConfig
}

func (t *tokenServer) GenerateToken(ctx context.Context, tr *proto.TokenRequest) (*proto.TokenString, error) {
	token, err := clientaccess.FormatToken(t.config.ControlConfig.Runtime.AgentToken, t.config.ControlConfig.Runtime.ServerCA)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new token: %w", err)
	}

	return &proto.TokenString{Token: token}, nil
}
