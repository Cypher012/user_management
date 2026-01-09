package token

import (
	"context"
	"time"

	"github.com/Cypher012/userauth/internal/security"
)

type TokenService struct {
	repo *TokenRepo
}

func NewTokenService(repo *TokenRepo) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) GetVerificationEmailToken(ctx context.Context, userId, email string) (rawToken string, err error) {
	rawToken, err = security.GenerateToken()
	if err != nil {
		return "", err
	}

	email_token_secret, err := security.GetEnv("EMAIL_TOKEN_SECRET")
	if err != nil {
		return "", err
	}

	hash := security.HashTokenKey(rawToken, email_token_secret)

	expires := time.Now().Add(1 * time.Hour)
	if err := s.repo.Create(ctx, userId, hash, VerifyEmailTokenType, expires); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *TokenService) VerifyEmailToken(ctx context.Context, rawToken string) error {
	email_token_secret, err := security.GetEnv("EMAIL_TOKEN_SECRET")
	if err != nil {
		return err
	}
	hash := security.HashTokenKey(rawToken, email_token_secret)
	token, err := s.repo.GetValidEmailToken(ctx, hash, VerifyEmailTokenType)
	if err != nil {
		return err
	}

	if err := s.repo.MarkEmailTokenUsed(ctx, token.ID.String()); err != nil {
		return err
	}

	return nil
}
