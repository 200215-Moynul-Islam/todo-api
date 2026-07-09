package services

import (
	"errors"
	"regexp"
	"strings"
	"todo-api/models"
	"todo-api/repositories"
	"todo-api/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrEmailExists        = errors.New("email already exists")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters long")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	RegisterUser(name, email, password string) (*models.User, error)
	LoginUser(email, password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (s *userService) RegisterUser(name, email, password string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" {
		return nil, ErrNameRequired
	}
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}

	if len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	// Check if email already exists
	existingEmail, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, ErrEmailExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) LoginUser(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return "", ErrEmailRequired
	}
	if password == "" {
		return "", ErrPasswordRequired
	}

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
