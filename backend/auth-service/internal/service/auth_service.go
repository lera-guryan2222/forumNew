package service

import (
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/usecase"
)

type AuthService interface {
	Register(req usecase.RegisterRequest) (*usecase.RegisterResponse, error)
	Login(req usecase.LoginRequest) (*usecase.LoginResponse, error)
	Refresh(req usecase.RefreshRequest) (*usecase.RefreshResponse, error)
}

type authService struct {
	uc usecase.AuthUsecase // Используем интерфейс, а не конкретную реализацию
}

// Принимаем интерфейс usecase.AuthUsecase
func NewAuthService(uc usecase.AuthUsecase) AuthService {
	return &authService{uc: uc}
}
func (s *authService) Register(req usecase.RegisterRequest) (*usecase.RegisterResponse, error) {
	return s.uc.Register(req)
}

func (s *authService) Login(req usecase.LoginRequest) (*usecase.LoginResponse, error) {
	return s.uc.Login(req)
}

func (s *authService) Refresh(req usecase.RefreshRequest) (*usecase.RefreshResponse, error) {
	return s.uc.Refresh(req)
}
