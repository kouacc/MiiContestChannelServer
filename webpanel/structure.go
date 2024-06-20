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
	Config
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Config struct {
	Username        string `xml:"username"`
	Password        string `xml:"password"`
	DatabaseAddress string `xml:"databaseAddress"`
	DatabaseName    string `xml:"databaseName"`
	Address         string `xml:"address"`
	AssetsPath      string `xml:"assetsPath"`
}
