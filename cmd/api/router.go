package main

import (
	"time"

	"github.com/Cypher012/userauth/internal/auth"
	"github.com/Cypher012/userauth/internal/http/v1/authhttp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool, jwt *auth.JWTAuth) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ok"))

	r.Route("/api", func(r chi.Router) {
		authhttp.RegisterAuth(r, pool, jwt)
	})

	return r
}
