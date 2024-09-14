package store

import (
	"context"

	"github.com/kapycluster/corpy/panel/model"
)

type DB struct {
}

func NewDB() (*DB, error) {
	return &DB{}, nil
}

func (d *DB) CreateControlPlane(ctx context.Context, cp *model.ControlPlane) error {
	return nil
}
