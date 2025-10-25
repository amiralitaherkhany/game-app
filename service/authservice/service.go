package authservice

import (
	"fmt"
	"gameapp/entity"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Config struct {
	SignKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

type Service struct {
	config Config
}

func New(c Config) *Service {
	return &Service{
		config: c,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	claims := AccessTokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   s.config.AccessSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := s.createNewJwtToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Service) CreateRefreshToken() (string, error) {
	claims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   s.config.RefreshSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := s.createNewJwtToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Service) createNewJwtToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.SignKey))
	return tokenString, err
}

func (s Service) ParseRefreshToken(bearerToken string) (*RefreshTokenClaims, error) {
	tokenString := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.SignKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return &RefreshTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return &RefreshTokenClaims{}, fmt.Errorf("invalid jwt token")
	}
}

func (s Service) ParseAccessToken(bearerToken string) (*AccessTokenClaims, error) {
	tokenString := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.SignKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return &AccessTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return &AccessTokenClaims{}, err
}
