package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type JWTAuth struct {
	tokenAuth *jwtauth.JWTAuth
}

type Claims struct {
	UserId string `json:"user_id"`
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"

	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 1 * time.Hour
)

func NewJWTAuth(secret string) *JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)
	return &JWTAuth{
		tokenAuth: tokenAuth,
	}
}

func (j *JWTAuth) GenerateToken(userID string) (
	atkToken string,
	rtkToken string,
	err error,
) {
	now := time.Now()

	atkClaims := map[string]any{
		"type":    AccessTokenType,
		"user_id": userID,
		"exp":     now.Add(AccessTokenTTL).Unix(),
	}

	_, atkToken, err = j.tokenAuth.Encode(atkClaims)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	rtkClaims := map[string]any{
		"type":    RefreshTokenType,
		"user_id": userID,
		"exp":     now.Add(RefreshTokenTTL).Unix(),
	}

	_, rtkToken, err = j.tokenAuth.Encode(rtkClaims)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return
}

func (j *JWTAuth) baseJWT(next http.Handler) http.Handler {
	verifier := jwtauth.Verifier(j.tokenAuth)
	authenticator := jwtauth.Authenticator(j.tokenAuth)
	return verifier(authenticator(next))
}

func (j *JWTAuth) AccessMiddleware(next http.Handler) http.Handler {
	return j.baseJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != AccessTokenType {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func (j *JWTAuth) RefreshMiddleware(next http.Handler) http.Handler {
	return j.baseJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != RefreshTokenType {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func (j *JWTAuth) FromContext(ctx context.Context) (Claims, error) {
	_, rawClaims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return Claims{}, err
	}
	user_id, ok := rawClaims["user_id"].(string)
	if !ok {
		return Claims{}, errors.New("user_id not found in claims")
	}

	tokenType, ok := rawClaims["type"].(string)
	if !ok || tokenType != AccessTokenType {
		return Claims{}, errors.New("invalid token type")
	}

	claims := Claims{
		UserId: user_id,
	}

	return claims, nil
}
