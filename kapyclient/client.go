package kapyclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type KapyClient struct {
	client *grpc.ClientConn
}

func NewKapyClient(address string) (*KapyClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating kapyclient: %w", err)
	}

	return &KapyClient{client: conn}, nil
}

func (k *KapyClient) Close() error {
	return k.client.Close()
}
