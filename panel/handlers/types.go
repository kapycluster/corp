package handlers

import (
	"context"

	"kapycluster.com/corp/panel/dns"
	"kapycluster.com/corp/panel/kube"
)

type KubeClient interface {
	CreateControlPlane(ctx context.Context, cp kube.ControlPlane) error
	UpdateControlPlane(ctx context.Context, cp kube.ControlPlane) error
	DeleteControlPlane(ctx context.Context, cp kube.ControlPlane) error
	WatchControlPlane(ctx context.Context, cp kube.ControlPlane) (<-chan bool, error)
	GetControlPlane(ctx context.Context, cp kube.ControlPlane) (*kube.ControlPlane, error)
	ListControlPlanes(ctx context.Context, userID string, regions []string) ([]*kube.ControlPlane, error)
	GetKubeconfig(ctx context.Context, cpID string, region string) ([]byte, error)
	ValidateControlPlane(cp kube.ControlPlane) error
	GetRegions() []string
}

type DBStore interface {
	CreateControlPlane(ctx context.Context, cp *kube.ControlPlane) error
	GetControlPlaneUser(ctx context.Context, cpID string) (string, error)
	GetUserControlPlanes(ctx context.Context, userID string) ([]*kube.ControlPlane, error)
	GetUserRegions(ctx context.Context, userID string) ([]string, error)
	GetControlPlane(ctx context.Context, cpID string) (*kube.ControlPlane, error)
	DeleteControlPlane(ctx context.Context, cpID string) error
}

type DNSClient interface {
	CreateDNSRecord(ctx context.Context, record dns.Record) error
	DeleteDNSRecord(ctx context.Context, recordID string) error
}
