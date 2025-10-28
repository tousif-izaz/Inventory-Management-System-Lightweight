package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/repository"
	"inventory-management/pkg/service/dto"
	"os"
	"time"
)

type IUserService interface {
	Login(username, password string) (string, error)
	SignUp(user dto.UserCreate) error
	GetUserByUsername(username string) (domain.User, error)
	UpdatePassword(username, oldPassword, newPassword string) error
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{userRepository}
}

func (service *UserService) Login(username, password string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	user, err := service.userRepository.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("no user found with the username: " + username)
	}

	// Check if user is active
	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Update last login timestamp
	_ = service.userRepository.UpdateLastLogin(user.UserID)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", errors.New("error signing the token: " + err.Error())
	}

	return tokenString, nil
}

func (service *UserService) SignUp(userCreate dto.UserCreate) error {
	err := validateUserCreate(userCreate)
	if err != nil {
		return err
	}

	user := userCreateToUser(userCreate)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error while creating password hash")
	}

	user.Password = string(hashedPassword)

	return service.userRepository.SignUp(user)
}

func validateUserCreate(u dto.UserCreate) error {
	if u.Username == "" {
		return errors.New("username can't be empty")
	}
	if u.Name == "" {
		return errors.New("name can't be empty")
	}
	if u.Password == "" {
		return errors.New("password can't be empty")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if u.Role == "" {
		return errors.New("role can't be empty")
	}
	// Validate role is one of the allowed values
	validRoles := map[string]bool{
		"admin":   true,
		"manager": true,
		"staff":   true,
	}
	if !validRoles[u.Role] {
		return errors.New("role must be one of: admin, manager, staff")
	}
	return nil
}

func userCreateToUser(userCreate dto.UserCreate) domain.User {
	return domain.User{
		Username: userCreate.Username,
		Name:     userCreate.Name,
		Email:    userCreate.Email,
		Password: userCreate.Password,
		Role:     userCreate.Role,
	}
}

func (service *UserService) GetUserByUsername(username string) (domain.User, error) {
	return service.userRepository.GetUserByUsername(username)
}

func (service *UserService) UpdatePassword(username, oldPassword, newPassword string) error {
	// Validate new password
	if newPassword == "" {
		return errors.New("new password can't be empty")
	}
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	// Get user
	user, err := service.userRepository.GetUserByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error while creating password hash")
	}

	// Update password
	return service.userRepository.UpdatePassword(user.UserID, string(hashedPassword))
}
