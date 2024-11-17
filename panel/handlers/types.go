package handlers

import (
	"context"

	"kapycluster.com/corp/panel/dns"
	"kapycluster.com/corp/panel/kube"
	"kapycluster.com/corp/panel/model"
)

type KubeClient interface {
	CreateControlPlane(ctx context.Context, cp kube.ControlPlane) error
	UpdateControlPlane(ctx context.Context, cp kube.ControlPlane) error
	DeleteControlPlane(ctx context.Context, cp kube.ControlPlane) error
	WatchControlPlane(ctx context.Context, cp kube.ControlPlane) (<-chan bool, error)
	GetControlPlane(ctx context.Context, cp kube.ControlPlane) (*kube.ControlPlane, error)
	ListControlPlanes(ctx context.Context, userID string) ([]*kube.ControlPlane, error)
}

type DBStore interface {
	CreateControlPlane(ctx context.Context, cp *model.ControlPlane) error
}

type DNSClient interface {
	CreateDNSRecord(ctx context.Context, record dns.Record) error
	DeleteDNSRecord(ctx context.Context, recordID string) error
}
