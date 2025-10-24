package authservice

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
}
