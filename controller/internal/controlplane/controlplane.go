package controlplane

import (
	"context"
	"fmt"

	"kapycluster.com/corp/controller/internal/controlplane/resources"
	"kapycluster.com/corp/controller/internal/scope"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Create is a helper function to create a control plane
func Create(ctx context.Context, client client.Client, scope *scope.ControlPlaneScope) error {
	l := log.FromContext(ctx, "control plane", scope.Name())

	l.Info("creating service for control plane")
	svc := resources.NewService(client, scope)
	if err := svc.Create(ctx); err != nil {
		return fmt.Errorf("control plane: failed to create service for %s: %w", scope.Name(), err)
	}

	l.Info("creating deployment for control plane")
	deploy := resources.NewDeployment(client, scope)
	if err := deploy.Create(ctx); err != nil {
		return fmt.Errorf("control plane: failed to create deployment for %s: %w", scope.Name(), err)
	}

	return nil
}
