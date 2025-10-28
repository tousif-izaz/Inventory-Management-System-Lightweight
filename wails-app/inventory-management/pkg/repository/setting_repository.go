package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"inventory-management/pkg/domain"

	"github.com/labstack/gommon/log"
)

type ISettingRepository interface {
	GetAllSettings() ([]domain.Setting, error)
	GetSettingByKey(key string) (domain.Setting, error)
	CreateSetting(setting domain.Setting) error
	UpdateSetting(key string, value string) error
	DeleteSetting(key string) error
}

type SettingRepository struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) ISettingRepository {
	return &SettingRepository{db}
}

func (repository *SettingRepository) GetAllSettings() ([]domain.Setting, error) {
	selectStatement := `
		SELECT SettingID, SettingKey, SettingValue, Description, CreatedAt, UpdatedAt
		FROM Settings
		ORDER BY SettingKey
	`

	rows, err := repository.db.Query(selectStatement)
	if err != nil {
		log.Errorf("error while fetching all settings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var settings []domain.Setting
	for rows.Next() {
		var setting domain.Setting
		var description sql.NullString

		err := rows.Scan(
			&setting.SettingID,
			&setting.SettingKey,
			&setting.SettingValue,
			&description,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			log.Errorf("error while scanning setting: %v", err)
			return nil, err
		}

		if description.Valid {
			setting.Description = &description.String
		}

		settings = append(settings, setting)
	}

	return settings, nil
}

func (repository *SettingRepository) GetSettingByKey(key string) (domain.Setting, error) {
	var setting domain.Setting
	var description sql.NullString

	selectStatement := `
		SELECT SettingID, SettingKey, SettingValue, Description, CreatedAt, UpdatedAt
		FROM Settings
		WHERE SettingKey = ?
	`

	err := repository.db.QueryRow(selectStatement, key).Scan(
		&setting.SettingID,
		&setting.SettingKey,
		&setting.SettingValue,
		&description,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.Setting{}, errors.New("setting not found")
	}

	if err != nil {
		log.Errorf("error while finding setting by key: %v", err)
		return domain.Setting{}, err
	}

	if description.Valid {
		setting.Description = &description.String
	}

	return setting, nil
}

func (repository *SettingRepository) CreateSetting(setting domain.Setting) error {
	insertStatement := `
		INSERT INTO Settings (SettingKey, SettingValue, Description)
		VALUES (?, ?, ?)
	`

	result, err := repository.db.Exec(
		insertStatement,
		setting.SettingKey,
		setting.SettingValue,
		setting.Description,
	)
	if err != nil {
		log.Errorf("error while creating setting: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Setting created successfully: %v", result))
	return nil
}

func (repository *SettingRepository) UpdateSetting(key string, value string) error {
	updateStatement := `
		UPDATE Settings
		SET SettingValue = ?
		WHERE SettingKey = ?
	`

	result, err := repository.db.Exec(updateStatement, value, key)
	if err != nil {
		log.Errorf("error while updating setting: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("setting not found")
	}

	log.Info(fmt.Sprintf("Setting updated successfully: %s = %s", key, value))
	return nil
}

func (repository *SettingRepository) DeleteSetting(key string) error {
	deleteStatement := `
		DELETE FROM Settings
		WHERE SettingKey = ?
	`

	result, err := repository.db.Exec(deleteStatement, key)
	if err != nil {
		log.Errorf("error while deleting setting: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("setting not found")
	}

	log.Info(fmt.Sprintf("Setting deleted successfully: %s", key))
	return nil
}
