package usecase

import (
	"errors"
	"time"

	"github.com/lera-guryan2222/forum/backend/auth-service/internal/entity"
	"github.com/lera-guryan2222/forum/backend/auth-service/internal/repository"
	"github.com/lera-guryan2222/forum/backend/auth-service/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(request LoginRequest) (*LoginResponse, error)
	Register(request RegisterRequest) (*RegisterResponse, error)
	Refresh(request RefreshRequest) (*RefreshResponse, error)
}

type authUsecase struct {
	userRepo     repository.UserRepository
	tokenRepo    repository.TokenRepository
	tokenManager auth.TokenManager
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	tokenManager auth.TokenManager,
) AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		tokenManager: tokenManager,
	}
}

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
		User         *entity.User `json:"user"`
	}

	RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterResponse struct {
		User         *entity.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
	}

	RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	RefreshResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func (uc *authUsecase) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.Save(user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (uc *authUsecase) Register(req RegisterRequest) (*RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("all fields are required")
	}

	existingUser, _ := uc.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingEmailUser, _ := uc.userRepo.FindByEmail(req.Email)
	if existingEmailUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.Save(user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authUsecase) Refresh(req RefreshRequest) (*RefreshResponse, error) {
	userID, expiresAt, err := uc.tokenRepo.Find(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(expiresAt) {
		return nil, errors.New("refresh token expired")
	}

	newAccessToken, err := uc.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, newExpiresAt, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.Delete(req.RefreshToken); err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.Save(userID, newRefreshToken, newExpiresAt); err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

var _ AuthUsecase = (*authUsecase)(nil)
