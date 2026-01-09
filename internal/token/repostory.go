package token

import (
	"context"
	"time"

	"github.com/Cypher012/userauth/internal/db/pgtypes"
	sqlc "github.com/Cypher012/userauth/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepo struct {
	q sqlc.Queries
}

func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{q: *sqlc.New(pool)}
}

func (r *TokenRepo) Create(ctx context.Context, userId, hash string, tokeType TokenType, expires time.Time) error {
	userUUID, err := pgtypes.ParseUUID(userId)
	if err != nil {
		return err
	}
	expiresPg, err := pgtypes.ParseTimestamp(expires)
	if err != nil {
		return err
	}

	return r.q.CreateEmailToken(ctx, sqlc.CreateEmailTokenParams{
		UserID:    *userUUID,
		TokenHash: hash,
		Type:      string(tokeType),
		ExpiresAt: *expiresPg,
	})
}

func (r *TokenRepo) GetValidEmailToken(ctx context.Context, hash string, tokeType TokenType) (sqlc.EmailToken, error) {
	return r.q.GetValidEmailToken(ctx, sqlc.GetValidEmailTokenParams{
		TokenHash: hash,
		Type:      string(tokeType),
	})
}

func (r *TokenRepo) MarkEmailTokenUsed(ctx context.Context, userId string) error {
	userUUID, err := pgtypes.ParseUUID(userId)
	if err != nil {
		return err
	}
	return r.q.MarkEmailTokenUsed(ctx, *userUUID)
}
