package domain

import "time"

// Setting represents a global application setting
type Setting struct {
	SettingID    int64     `json:"setting_id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
