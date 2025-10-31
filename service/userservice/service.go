package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
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

func (s Service) GetProfile(req dto.GetProfileRequest) (dto.GetProfileResponse, error) {
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.GetProfileResponse{}, richerror.New("").WithErr(err)
	}

	return dto.GetProfileResponse{
		Name: user.Name,
	}, nil
}

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{},
			richerror.New("userservice.Login").WithErr(err)
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return dto.LoginResponse{}, fmt.Errorf("phone number or password isn't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken()
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.LoginResponse{
		User: dto.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Tokens: dto.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// hash password
	hash, err := HashPassword(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("can't hash your password")
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
		return dto.RegisterResponse{}, richerror.New("").WithErr(err)
	}

	return dto.RegisterResponse{
		User: dto.UserInfo{
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
