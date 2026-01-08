package pgtypes

import "github.com/jackc/pgx/v5/pgtype"

func ParseUUID(id string) (*pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(id); err != nil {
		return nil, err
	}
	return &u, nil
}
