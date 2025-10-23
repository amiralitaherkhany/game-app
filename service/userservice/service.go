package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, error)
	GetUserByID(userID uint) (*entity.User, error)
}

type Service struct {
	signKey string
	repo    Repository
}

type GetProfileRequest struct {
	UserID uint `json:"user_id"`
}

type GetProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) GetProfile(req GetProfileRequest) (GetProfileResponse, error) {
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return GetProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if user == nil {
		return GetProfileResponse{}, fmt.Errorf("user not found")
	}

	return GetProfileResponse{
		Name: user.Name,
	}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if user == nil {
		return LoginResponse{}, fmt.Errorf("phone number or password isn't correct")
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return LoginResponse{}, fmt.Errorf("phone number or password isn't correct")
	}

	return LoginResponse{}, nil
}

func New(repo Repository, signKey string) *Service {
	return &Service{repo: repo, signKey: signKey}
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}
type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check the uniqueness of phone number
	isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	if !isUnique {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// TODO - we must verify phone number by verification code

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be 3 or greater")
	}

	// TODO - check password with regex pattern
	// validate password
	if len(req.Password) < 8 || len(req.Password) > 20 {
		return RegisterResponse{}, fmt.Errorf("password length should be 8-20 byte")
	}

	// hash password
	hash, err := HashPassword(req.Password)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("can't hash your password")
	}

	// register user
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hash,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	jwtToken, err := createNewJwtToken(createdUser.ID, s.signKey)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return RegisterResponse{
		AccessToken: jwtToken,
	}, nil
}

func HashPassword(password string) (string, error) {
	// password length can't be greater than 72 bytes
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type MyCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func createNewJwtToken(userID uint, signKey string) (string, error) {
	claims := MyCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signKey))

	return tokenString, err
}
