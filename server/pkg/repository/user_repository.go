package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

type IUserRepository interface {
	GetUserByUsername(username string) (domain.User, error)
	SignUp(user domain.User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db}
}

func (repository *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User

	selectStatement := "SELECT id, username, password, role FROM users WHERE username = ?"
	err := repository.db.QueryRow(selectStatement, username).Scan(&user.Id, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("error while finding user")
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) SignUp(user domain.User) error {
	insertStatement := "INSERT INTO users(username, password, role) VALUES (?, ?, ?)"

	result, err := repository.db.Exec(insertStatement, user.Username, user.Password, user.Role)
	if err != nil {
		log.Errorf("error while adding new user: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("User added successfully: %v", result))
	return nil
}
