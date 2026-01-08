package authhttp

import (
	"net/http"

	"github.com/Cypher012/userauth/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuth(r chi.Router, pool *pgxpool.Pool, jwt *auth.JWTAuth) {
	repo := auth.NewAuthRepository(pool)
	service := auth.NewAuthService(repo)
	handler := NewAuthHandler(service, jwt)

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/signup", handler.SignUpHandler)
		r.Post("/login", handler.LoginHandler)
		r.Post("/verify-email", handler.VerifyEmailHandler)
		r.Post("/send-verification-email", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/forget-password", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {})
		r.With(jwt.RefreshMiddleware).Post("/refresh", func(w http.ResponseWriter, r *http.Request) {})
	})
}
