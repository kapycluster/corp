package store

import (
	"context"

	"kapycluster.com/corp/panel/model"
)

type DB struct {
}

func NewDB() (*DB, error) {
	return &DB{}, nil
}

func (d *DB) CreateControlPlane(ctx context.Context, cp *model.ControlPlane) error {
	return nil
}
