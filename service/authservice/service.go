package authservice

import (
	"fmt"
	"gameapp/entity"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func New(signKey, accessSubject, refreshSubject string, accessExpirationTime, refreshExpirationTime time.Duration) *Service {
	return &Service{
		signKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	claims := AccessTokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   s.accessSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpirationTime)),
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
			Subject:   s.refreshSubject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpirationTime)),
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
	tokenString, err := token.SignedString([]byte(s.signKey))
	return tokenString, err
}

func (s Service) ParseRefreshToken(bearerToken string) (*RefreshTokenClaims, error) {
	tokenString := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.signKey), nil
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
		return []byte(s.signKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return &AccessTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return &AccessTokenClaims{}, err
}
