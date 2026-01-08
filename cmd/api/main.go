package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Cypher012/userauth/internal/auth"
	"github.com/Cypher012/userauth/internal/common"
	"github.com/Cypher012/userauth/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret, err := common.GetEnv("JWT_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	jwt := auth.NewJWTAuth(secret)

	dbPool := db.NewDB(ctx)
	defer dbPool.Close()

	r := NewRouter(dbPool, jwt)

	port, err := common.GetEnv("PORT")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server running on :%v", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
