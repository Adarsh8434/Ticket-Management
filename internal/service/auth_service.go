package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"ticket-system/internal/model"
	"ticket-system/internal/repository"
)

type AuthService struct {
	UserRepository *repository.UserRepository
	JWTSecret      string
}

func NewAuthService(
	userRepository *repository.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		JWTSecret:      jwtSecret,
	}
}

// Register creates a new user.
func (s *AuthService) Register(
	email string,
	password string,
) (*model.User, error) {

	// Check if user already exists
	_, err := s.UserRepository.FindByEmail(email)

	if err == nil {
		return nil, errors.New("email already registered")
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	// Save user
	return s.UserRepository.CreateUser(
		email,
		string(hashedPassword),
	)
}

// Login authenticates a user and generates JWT.
func (s *AuthService) Login(
	email string,
	password string,
) (string, error) {

	// Find user
	user, err := s.UserRepository.FindByEmail(email)

	if err != nil {
		return "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", err
	}

	// JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create JWT
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign JWT
	tokenString, err := token.SignedString(
		[]byte(s.JWTSecret),
	)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
