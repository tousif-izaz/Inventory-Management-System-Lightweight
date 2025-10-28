package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"inventory-management/pkg/domain"
	"time"

	"github.com/labstack/gommon/log"
)

type IUserRepository interface {
	GetUserByUsername(username string) (domain.User, error)
	GetUserByID(userID int64) (domain.User, error)
	SignUp(user domain.User) error
	UpdateLastLogin(userID int64) error
	UpdatePassword(userID int64, newPasswordHash string) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db}
}

func (repository *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	var email sql.NullString
	var lastLogin sql.NullTime

	selectStatement := `
		SELECT UserID, Username, Name, Email, PasswordHash, Role, IsActive, CreatedAt, UpdatedAt, LastLogin
		FROM Users
		WHERE Username = ? AND IsActive = 1
	`
	err := repository.db.QueryRow(selectStatement, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Name,
		&email,
		&user.Password, // PasswordHash maps to Password field in domain
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found or inactive")
	}

	if err != nil {
		log.Errorf("error while finding user: %v", err)
		return domain.User{}, err
	}

	// Handle nullable fields
	if email.Valid {
		user.Email = &email.String
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

func (repository *UserRepository) GetUserByID(userID int64) (domain.User, error) {
	var user domain.User
	var email sql.NullString
	var lastLogin sql.NullTime

	selectStatement := `
		SELECT UserID, Username, Name, Email, PasswordHash, Role, IsActive, CreatedAt, UpdatedAt, LastLogin
		FROM Users
		WHERE UserID = ?
	`
	err := repository.db.QueryRow(selectStatement, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Name,
		&email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}

	if err != nil {
		log.Errorf("error while finding user by ID: %v", err)
		return domain.User{}, err
	}

	if email.Valid {
		user.Email = &email.String
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

func (repository *UserRepository) SignUp(user domain.User) error {
	// Extended schema requires: Username, Name, Email (optional), PasswordHash, Role
	insertStatement := `
		INSERT INTO Users (Username, Name, Email, PasswordHash, Role, IsActive)
		VALUES (?, ?, ?, ?, ?, 1)
	`

	result, err := repository.db.Exec(
		insertStatement,
		user.Username,
		user.Name,
		user.Email,
		user.Password, // Password field contains the hash
		user.Role,
	)
	if err != nil {
		log.Errorf("error while adding new user: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("User added successfully: %v", result))
	return nil
}

func (repository *UserRepository) UpdateLastLogin(userID int64) error {
	updateStatement := `
		UPDATE Users
		SET LastLogin = ?
		WHERE UserID = ?
	`

	_, err := repository.db.Exec(updateStatement, time.Now(), userID)
	if err != nil {
		log.Errorf("error while updating last login: %v", err)
		return err
	}

	return nil
}

func (repository *UserRepository) UpdatePassword(userID int64, newPasswordHash string) error {
	updateStatement := `
		UPDATE Users
		SET PasswordHash = ?
		WHERE UserID = ?
	`

	_, err := repository.db.Exec(updateStatement, newPasswordHash, userID)
	if err != nil {
		log.Errorf("error while updating password: %v", err)
		return err
	}

	return nil
}
