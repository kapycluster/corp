package handlers

import (
	"context"

	"github.com/kapycluster/corpy/panel/kube"
	"github.com/kapycluster/corpy/panel/model"
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
