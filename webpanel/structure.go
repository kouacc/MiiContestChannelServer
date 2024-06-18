package webpanel

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WebPanel struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
	Salt []byte
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
