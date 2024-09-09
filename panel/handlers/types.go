package handlers

import (
	"context"

	kapyv1 "github.com/kapycluster/corpy/controller/api/v1"
)

type ControlPlaneClient interface {
	CreateControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error
	UpdateControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error
	DeleteControlPlane(ctx context.Context, cp kapyv1.ControlPlane) error
	WatchControlPlane(ctx context.Context, cp kapyv1.ControlPlane) (<-chan bool, error)
}
