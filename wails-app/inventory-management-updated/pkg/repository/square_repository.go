package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"ims-intro/pkg/domain"

	"github.com/labstack/gommon/log"
)

type ISquareRepository interface {
	// Token management
	SaveOAuthToken(token domain.SquareOAuthToken) error
	GetOAuthTokenByMerchant(merchantID string) (domain.SquareOAuthToken, error)
	UpdateOAuthToken(merchantID string, token domain.SquareOAuthToken) error
	DeleteOAuthToken(merchantID string) error
	GetAllMerchantIDs() ([]string, error)

	// Config management (uses Settings table)
	SaveConfig(config domain.SquareConfig) error
	GetConfig() (domain.SquareConfig, error)
	SaveMerchantID(merchantID string) error
	GetMerchantID() (string, error)
}

type SquareRepository struct {
	db *sql.DB
}

func NewSquareRepository(db *sql.DB) ISquareRepository {
	return &SquareRepository{db}
}

// SaveOAuthToken saves a new OAuth token
func (r *SquareRepository) SaveOAuthToken(token domain.SquareOAuthToken) error {
	insertStatement := `
		INSERT INTO SquareOAuthTokens (MerchantID, AccessToken, RefreshToken, ExpiresAt, TokenType, Scopes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		insertStatement,
		token.MerchantID,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
		token.TokenType,
		token.Scopes,
	)
	if err != nil {
		log.Errorf("error while saving Square OAuth token: %v", err)
		return err
	}

	log.Info("Square OAuth token saved successfully")
	return nil
}

// GetOAuthTokenByMerchant retrieves an OAuth token by merchant ID
func (r *SquareRepository) GetOAuthTokenByMerchant(merchantID string) (domain.SquareOAuthToken, error) {
	var token domain.SquareOAuthToken
	var expiresAt sql.NullTime
	var refreshToken sql.NullString

	selectStatement := `
		SELECT TokenID, MerchantID, AccessToken, RefreshToken, ExpiresAt, TokenType, Scopes, CreatedAt, UpdatedAt
		FROM SquareOAuthTokens
		WHERE MerchantID = ?
	`

	err := r.db.QueryRow(selectStatement, merchantID).Scan(
		&token.TokenID,
		&token.MerchantID,
		&token.AccessToken,
		&refreshToken,
		&expiresAt,
		&token.TokenType,
		&token.Scopes,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.SquareOAuthToken{}, errors.New("token not found")
	}

	if err != nil {
		log.Errorf("error while getting Square OAuth token: %v", err)
		return domain.SquareOAuthToken{}, err
	}

	if refreshToken.Valid {
		token.RefreshToken = refreshToken.String
	}

	if expiresAt.Valid {
		token.ExpiresAt = &expiresAt.Time
	}

	return token, nil
}

// UpdateOAuthToken updates an existing OAuth token
func (r *SquareRepository) UpdateOAuthToken(merchantID string, token domain.SquareOAuthToken) error {
	updateStatement := `
		UPDATE SquareOAuthTokens
		SET AccessToken = ?, RefreshToken = ?, ExpiresAt = ?, TokenType = ?, Scopes = ?
		WHERE MerchantID = ?
	`

	result, err := r.db.Exec(
		updateStatement,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
		token.TokenType,
		token.Scopes,
		merchantID,
	)
	if err != nil {
		log.Errorf("error while updating Square OAuth token: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("token not found")
	}

	log.Info(fmt.Sprintf("Square OAuth token updated successfully for merchant: %s", merchantID))
	return nil
}

// DeleteOAuthToken deletes an OAuth token
func (r *SquareRepository) DeleteOAuthToken(merchantID string) error {
	deleteStatement := `
		DELETE FROM SquareOAuthTokens
		WHERE MerchantID = ?
	`

	result, err := r.db.Exec(deleteStatement, merchantID)
	if err != nil {
		log.Errorf("error while deleting Square OAuth token: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("token not found")
	}

	log.Info(fmt.Sprintf("Square OAuth token deleted successfully for merchant: %s", merchantID))
	return nil
}

// SaveConfig saves Square configuration to Settings table
func (r *SquareRepository) SaveConfig(config domain.SquareConfig) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Errorf("error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Save each config value as a setting
	settings := map[string]string{
		"square_application_id":     config.ApplicationID,
		"square_application_secret": config.ApplicationSecret,
		"square_environment":        config.Environment,
		"square_redirect_uri":       config.RedirectURI,
	}

	for key, value := range settings {
		// Try to update first
		updateStmt := `UPDATE Settings SET SettingValue = ? WHERE SettingKey = ?`
		result, err := tx.Exec(updateStmt, value, key)
		if err != nil {
			log.Errorf("error updating setting %s: %v", key, err)
			return err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			// If no rows updated, insert new
			insertStmt := `INSERT INTO Settings (SettingKey, SettingValue, Description) VALUES (?, ?, ?)`
			description := fmt.Sprintf("Square %s configuration", key[7:]) // Remove "square_" prefix
			_, err = tx.Exec(insertStmt, key, value, description)
			if err != nil {
				log.Errorf("error inserting setting %s: %v", key, err)
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("error committing transaction: %v", err)
		return err
	}

	log.Info("Square configuration saved successfully")
	return nil
}

// GetConfig retrieves Square configuration from Settings table
func (r *SquareRepository) GetConfig() (domain.SquareConfig, error) {
	var config domain.SquareConfig

	query := `
		SELECT SettingKey, SettingValue
		FROM Settings
		WHERE SettingKey IN ('square_application_id', 'square_application_secret', 'square_environment', 'square_redirect_uri')
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return domain.SquareConfig{}, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			log.Errorf("error scanning setting: %v", err)
			return domain.SquareConfig{}, err
		}
		settings[key] = value
	}

	config.ApplicationID = settings["square_application_id"]
	config.ApplicationSecret = settings["square_application_secret"]
	config.Environment = settings["square_environment"]
	config.RedirectURI = settings["square_redirect_uri"]

	// Validate that we have required config
	if config.ApplicationID == "" || config.ApplicationSecret == "" {
		return domain.SquareConfig{}, errors.New("Square configuration not found or incomplete")
	}

	return config, nil
}

// SaveMerchantID saves the merchant ID to settings
func (r *SquareRepository) SaveMerchantID(merchantID string) error {
	// Try to update first
	updateStmt := `UPDATE Settings SET SettingValue = ? WHERE SettingKey = ?`
	result, err := r.db.Exec(updateStmt, merchantID, "square_merchant_id")
	if err != nil {
		log.Errorf("error updating square_merchant_id: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// If no rows updated, insert new
		insertStmt := `INSERT INTO Settings (SettingKey, SettingValue, Description) VALUES (?, ?, ?)`
		_, err = r.db.Exec(insertStmt, "square_merchant_id", merchantID, "Square merchant ID from OAuth")
		if err != nil {
			log.Errorf("error inserting square_merchant_id: %v", err)
			return err
		}
	}

	log.Info(fmt.Sprintf("Saved merchant ID: %s", merchantID))
	return nil
}

// GetMerchantID retrieves the merchant ID from settings
func (r *SquareRepository) GetMerchantID() (string, error) {
	var merchantID string
	query := `SELECT SettingValue FROM Settings WHERE SettingKey = 'square_merchant_id'`
	err := r.db.QueryRow(query).Scan(&merchantID)
	if err == sql.ErrNoRows {
		return "", errors.New("merchant ID not found")
	}
	if err != nil {
		log.Errorf("error getting merchant ID: %v", err)
		return "", err
	}
	return merchantID, nil
}

// GetAllMerchantIDs retrieves all merchant IDs from SquareOAuthTokens table
func (r *SquareRepository) GetAllMerchantIDs() ([]string, error) {
	query := `SELECT MerchantID FROM SquareOAuthTokens`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Errorf("error querying merchant IDs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var merchantIDs []string
	for rows.Next() {
		var merchantID string
		if err := rows.Scan(&merchantID); err != nil {
			log.Errorf("error scanning merchant ID: %v", err)
			continue
		}
		merchantIDs = append(merchantIDs, merchantID)
	}

	return merchantIDs, nil
}
