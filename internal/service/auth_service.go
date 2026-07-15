package service

import (
	"errors"

	"github.com/you/sharing-vision-backend-v2/internal/auth"
	"github.com/you/sharing-vision-backend-v2/internal/model"
)

type UserRepo interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	CountAdmins() (int, error)
}

type AuthService struct {
	userRepo UserRepo
}

func NewAuthService(userRepo UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.LoginResponse, error) {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !auth.ComparePassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: *user}, nil
}

func (s *AuthService) BootstrapAdmin() error {
	hash, err := auth.HashPassword("admin123")
	if err != nil {
		return err
	}
	u, err := s.userRepo.FindByEmail("admin@sharingvision.id")
	if err != nil {
		return err
	}
	if u == nil {
		return s.userRepo.Create(&model.User{
			Email:        "admin@sharingvision.id",
			PasswordHash: hash,
			Role:         "admin",
		})
	}
	// Admin sudah ada, skip bootstrap
	return nil
}
