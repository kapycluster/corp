package store

import (
	"context"
	"fmt"

	"kapycluster.com/corp/panel/model"
)

func (db *DB) GetInvite(ctx context.Context, inviteID string) (*model.Invite, error) {
	var invite model.Invite
	var usedInt int
	if err := db.QueryRowContext(ctx,
		"SELECT id, used FROM invites WHERE id = ?",
		inviteID,
	).Scan(&invite.ID, &usedInt); err != nil {
		return nil, fmt.Errorf("querying invite: %w", err)
	}

	invite.Used = usedInt == 1
	return &invite, nil
}

func (db *DB) UseInvite(ctx context.Context, inviteID string) error {
	if _, err := db.ExecContext(ctx,
		"UPDATE invites SET used = 1 WHERE id = ?",
		inviteID,
	); err != nil {
		return fmt.Errorf("updating invite: %w", err)
	}
	return nil
}
