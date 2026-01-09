package pgtypes

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(id string) (*pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(id); err != nil {
		return nil, err
	}
	return &u, nil
}

func ParseTimestamp(timestamp time.Time) (*pgtype.Timestamp, error) {
	var t pgtype.Timestamp
	if err := t.Scan(timestamp); err != nil {
		return nil, err
	}
	return &t, nil
}
