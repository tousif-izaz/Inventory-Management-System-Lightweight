package controller

import (
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"net/http"
)

type SettingController struct {
	settingService service.ISettingService
}

func NewSettingController(settingService service.ISettingService) *SettingController {
	return &SettingController{settingService}
}

func (controller *SettingController) RegisterSettingRoutes(e *echo.Echo) {
	settingsGroup := e.Group("/settings")
	settingsGroup.Use(middleware.AuthMiddleware)

	settingsGroup.GET("", controller.GetAllSettings)
	settingsGroup.GET("/:key", controller.GetSettingByKey)
	settingsGroup.POST("", controller.CreateSetting)
	settingsGroup.PUT("/:key", controller.UpdateSetting)
	settingsGroup.DELETE("/:key", controller.DeleteSetting)
}

func (controller *SettingController) GetAllSettings(c echo.Context) error {
	settings, err := controller.settingService.GetAllSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, settings)
}

func (controller *SettingController) GetSettingByKey(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no setting key specified"))
	}

	setting, err := controller.settingService.GetSettingByKey(key)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, setting)
}

func (controller *SettingController) CreateSetting(c echo.Context) error {
	var setting domain.Setting

	err := c.Bind(&setting)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the setting structure"))
	}

	err = controller.settingService.CreateSetting(setting)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *SettingController) UpdateSetting(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no setting key specified"))
	}

	var updateRequest struct {
		SettingValue string `json:"setting_value"`
	}

	err := c.Bind(&updateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data"))
	}

	err = controller.settingService.UpdateSetting(key, updateRequest.SettingValue)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (controller *SettingController) DeleteSetting(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: no setting key specified"))
	}

	err := controller.settingService.DeleteSetting(key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusNoContent)
}
