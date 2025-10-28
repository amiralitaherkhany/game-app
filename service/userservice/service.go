package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"
	"gameapp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthService interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken() (string, error)
}

type Service struct {
	auth AuthService
	repo Repository
}

func New(repo Repository, auth AuthService) *Service {
	return &Service{repo: repo, auth: auth}
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
		return GetProfileResponse{}, richerror.New("").WithErr(err)
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
	User   UserInfo `json:"user"`
	Tokens Tokens   `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{},
			richerror.New("userservice.Login").WithErr(err)
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return LoginResponse{}, fmt.Errorf("phone number or password isn't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{
		User: UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RegisterResponse struct {
	User UserInfo `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check the uniqueness of phone number
	isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber)
	if err != nil {
		return RegisterResponse{}, richerror.New("").WithErr(err)
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
		return RegisterResponse{}, richerror.New("").WithErr(err)
	}

	return RegisterResponse{
		User: UserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			PhoneNumber: createdUser.PhoneNumber,
		},
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
