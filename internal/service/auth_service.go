package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cachewraith/task-list-api-golang/internal/domain"
	"github.com/cachewraith/task-list-api-golang/internal/repository"
	"github.com/cachewraith/task-list-api-golang/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// AuthService defines the interface for authentication
type AuthService interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input domain.LoginInput) (string, *domain.User, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// authService implements AuthService
type authService struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// ensure authService implements AuthService
var _ AuthService = (*authService)(nil)

// Register creates a new user account
func (s *authService) Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error) {
	logger.Info("AUTH_SERVICE", "Register", "Attempting to register user with email: "+input.Email)

	if err := input.Validate(); err != nil {
		logger.Warn("AUTH_SERVICE", "Register", "Validation failed: "+err.Error())
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, input.Email)
	if existingUser != nil {
		logger.Warn("AUTH_SERVICE", "Register", "Duplicate email registration attempt: "+input.Email)
		return nil, domain.ErrDuplicate("user with this email")
	}

	// Create new user
	user := domain.NewUser(input)

	// Hash password
	if err := user.HashPassword(); err != nil {
		logger.Error("AUTH_SERVICE", "Register", "Failed to hash password: "+err.Error())
		return nil, domain.ErrInvalidInput("failed to hash password")
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("AUTH_SERVICE", "Register", "Failed to create user in database: "+err.Error())
		return nil, err
	}

	logger.Info("AUTH_SERVICE", "Register", "Successfully registered user: "+user.ID.String())
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, input domain.LoginInput) (string, *domain.User, error) {
	logger.Info("AUTH_SERVICE", "Login", "Login attempt for email: "+input.Email)

	if err := input.Validate(); err != nil {
		logger.Warn("AUTH_SERVICE", "Login", "Validation failed: "+err.Error())
		return "", nil, err
	}

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if domain.IsNotFound(err) {
			logger.Warn("AUTH_SERVICE", "Login", "Login failed - user not found: "+input.Email)
			return "", nil, domain.ErrValidationFailed("invalid email or password")
		}
		logger.Error("AUTH_SERVICE", "Login", "Database error during login: "+err.Error())
		return "", nil, err
	}

	// Verify password
	if !user.CheckPassword(input.Password) {
		logger.Warn("AUTH_SERVICE", "Login", "Login failed - invalid password for user: "+input.Email)
		return "", nil, domain.ErrValidationFailed("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.Error("AUTH_SERVICE", "Login", "Failed to generate token: "+err.Error())
		return "", nil, err
	}

	logger.Info("AUTH_SERVICE", "Login", "Successful login for user: "+user.ID.String())
	return token, user, nil
}

// generateToken creates a new JWT token for a user
func (s *authService) generateToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "todo-api",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(tokenString string) (*JWTClaims, error) {
	logger.Debug("AUTH_SERVICE", "ValidateToken", "Validating token...")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("AUTH_SERVICE", "ValidateToken", "Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		logger.Warn("AUTH_SERVICE", "ValidateToken", "Token validation failed: "+err.Error())
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		logger.Debug("AUTH_SERVICE", "ValidateToken", "Token validated successfully for user: "+claims.UserID.String())
		return claims, nil
	}

	logger.Warn("AUTH_SERVICE", "ValidateToken", "Invalid token claims")
	return nil, fmt.Errorf("invalid token claims")
}
