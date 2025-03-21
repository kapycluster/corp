package store

import (
	"context"
	"fmt"

	"kapycluster.com/corp/panel/kube"
)

func (db *DB) CreateControlPlane(ctx context.Context, cp *kube.ControlPlane) error {
	if cp == nil {
		return fmt.Errorf("control plane cannot be nil")
	}

	_, err := db.ExecContext(ctx,
		"INSERT INTO control_planes (id, name, user_id, region) VALUES (?, ?, ?, ?)",
		cp.ID, cp.Name, cp.UserID, cp.Region,
	)
	if err != nil {
		return fmt.Errorf("creating control plane: %w", err)
	}
	return nil
}

func (db *DB) GetUserControlPlanes(ctx context.Context, userID string) ([]*kube.ControlPlane, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT id, name, user_id, region FROM control_planes WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying control planes: %w", err)
	}
	defer rows.Close()

	var controlPlanes []*kube.ControlPlane
	for rows.Next() {
		cp := &kube.ControlPlane{}
		if err := rows.Scan(&cp.ID, &cp.Name, &cp.UserID, &cp.Region); err != nil {
			return nil, fmt.Errorf("scanning control plane: %w", err)
		}
		controlPlanes = append(controlPlanes, cp)
	}
	return controlPlanes, nil
}

func (db *DB) DeleteControlPlane(ctx context.Context, cpID string) error {
	_, err := db.ExecContext(ctx,
		"DELETE FROM control_planes WHERE id = ?",
		cpID,
	)
	if err != nil {
		return fmt.Errorf("deleting control plane: %w", err)
	}
	return nil
}

func (db *DB) GetControlPlane(ctx context.Context, cpID string) (*kube.ControlPlane, error) {
	cp := &kube.ControlPlane{}
	if err := db.QueryRowContext(ctx,
		"SELECT id, name, user_id, region FROM control_planes WHERE id = ?",
		cpID,
	).Scan(&cp.ID, &cp.Name, &cp.UserID, &cp.Region); err != nil {
		return nil, fmt.Errorf("querying control plane: %w", err)
	}
	return cp, nil
}

func (db *DB) GetControlPlaneUser(ctx context.Context, cpID string) (string, error) {
	var userID string
	if err := db.QueryRowContext(ctx,
		"SELECT user_id FROM control_planes WHERE id = ?",
		cpID,
	).Scan(&userID); err != nil {
		return "", fmt.Errorf("querying control plane user: %w", err)
	}
	return userID, nil
}

func (db *DB) GetUserRegions(ctx context.Context, userID string) ([]string, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT DISTINCT region FROM control_planes WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying regions: %w", err)
	}
	defer rows.Close()

	var regions []string
	for rows.Next() {
		var region string
		if err := rows.Scan(&region); err != nil {
			return nil, fmt.Errorf("scanning region: %w", err)
		}
		regions = append(regions, region)
	}
	return regions, nil
}
