package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/internal/repository"
	"musicapp/backend/pkg/apierror"
	"musicapp/backend/pkg/token"
)

type RegisterInput struct {
	Email       string `json:"email"        binding:"required,email"`
	Username    string `json:"username"     binding:"required,min=3,max=50"`
	Password    string `json:"password"     binding:"required,min=8"`
	DisplayName string `json:"display_name"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateMeInput struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

type AuthService struct {
	userRepo      repository.UserRepository
	tokenSvc      *token.Service
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	db            *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenSvc *token.Service,
	tokenExpiry time.Duration,
	refreshExpiry time.Duration,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		tokenSvc:      tokenSvc,
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
		db:            db,
	}
}

func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, apierror.New(409, "EMAIL_TAKEN", "an account with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("checking email availability: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: string(hash),
		DisplayName:  input.DisplayName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return s.issueTokens(user)
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.New(401, "INVALID_CREDENTIALS", "invalid email or password")
		}
		return nil, fmt.Errorf("finding user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, apierror.New(401, "INVALID_CREDENTIALS", "invalid email or password")
	}

	return s.issueTokens(user)
}

func (s *AuthService) RefreshTokens(rawToken string) (*AuthResponse, error) {
	hashed := hashToken(rawToken)

	var rt models.RefreshToken
	if err := s.db.Where("token_hash = ?", hashed).First(&rt).Error; err != nil {
		return nil, apierror.New(401, "INVALID_REFRESH_TOKEN", "refresh token is invalid or expired")
	}

	if time.Now().After(rt.ExpiresAt) {
		s.db.Delete(&rt)
		return nil, apierror.New(401, "REFRESH_TOKEN_EXPIRED", "refresh token has expired, please log in again")
	}

	s.db.Delete(&rt)

	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("finding user for refresh: %w", err)
	}

	return s.issueTokens(user)
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) GetMe(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) UpdateMe(userID uuid.UUID, input UpdateMeInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if input.DisplayName != "" {
		user.DisplayName = input.DisplayName
	}
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("updating profile: %w", err)
	}
	return user, nil
}

func (s *AuthService) issueTokens(user *models.User) (*AuthResponse, error) {
	accessToken, err := s.tokenSvc.Generate(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	rawRefresh, err := generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.refreshExpiry),
	}
	if err := s.db.Create(rt).Error; err != nil {
		return nil, fmt.Errorf("saving refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.tokenExpiry.Seconds()),
		RefreshToken: rawRefresh,
		User:         user,
	}, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
