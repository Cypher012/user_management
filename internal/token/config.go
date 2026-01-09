package token

import "github.com/jackc/pgx/v5/pgxpool"

func TokenConfig(pool *pgxpool.Pool) *TokenService {
	repo := NewTokenRepo(pool)
	service := NewTokenService(repo)

	return service
}
