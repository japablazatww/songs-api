package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppAuth struct {
	ID           uint64    `json:"id"`
	AppName      string    `json:"app_name"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Claims struct {
	jwt.RegisteredClaims
	ClientID string `json:"client_id"`
	AppName  string `json:"app_name"`
}
