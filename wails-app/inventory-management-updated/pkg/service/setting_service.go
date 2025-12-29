package service

import (
	"errors"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
)

type ISettingService interface {
	GetAllSettings() ([]domain.Setting, error)
	GetSettingByKey(key string) (domain.Setting, error)
	GetSettingValue(key string) (string, error)
	CreateSetting(setting domain.Setting) error
	UpdateSetting(key string, value string) error
	DeleteSetting(key string) error
}

type SettingService struct {
	settingRepository repository.ISettingRepository
}

func NewSettingService(settingRepository repository.ISettingRepository) ISettingService {
	return &SettingService{settingRepository}
}

func (service *SettingService) GetAllSettings() ([]domain.Setting, error) {
	return service.settingRepository.GetAllSettings()
}

func (service *SettingService) GetSettingByKey(key string) (domain.Setting, error) {
	if key == "" {
		return domain.Setting{}, errors.New("setting key cannot be empty")
	}
	return service.settingRepository.GetSettingByKey(key)
}

func (service *SettingService) GetSettingValue(key string) (string, error) {
	setting, err := service.GetSettingByKey(key)
	if err != nil {
		return "", err
	}
	return setting.SettingValue, nil
}

func (service *SettingService) CreateSetting(setting domain.Setting) error {
	if err := validateSetting(setting); err != nil {
		return err
	}
	return service.settingRepository.CreateSetting(setting)
}

func (service *SettingService) UpdateSetting(key string, value string) error {
	if key == "" {
		return errors.New("setting key cannot be empty")
	}
	if value == "" {
		return errors.New("setting value cannot be empty")
	}
	return service.settingRepository.UpdateSetting(key, value)
}

func (service *SettingService) DeleteSetting(key string) error {
	if key == "" {
		return errors.New("setting key cannot be empty")
	}
	return service.settingRepository.DeleteSetting(key)
}

func validateSetting(setting domain.Setting) error {
	if setting.SettingKey == "" {
		return errors.New("setting key cannot be empty")
	}
	if setting.SettingValue == "" {
		return errors.New("setting value cannot be empty")
	}
	return nil
}
