package controller

import (
	"github.com/labstack/echo/v4"
	"inventory-management/pkg/backup"
	"inventory-management/pkg/controller/response"
	"net/http"
)

type BackupController struct {
	backupService *backup.BackupService
}

func NewBackupController(backupService *backup.BackupService) *BackupController {
	return &BackupController{backupService: backupService}
}

func (controller *BackupController) RegisterBackupRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	// All backup routes require authentication
	e.POST("/backup", controller.CreateBackup, authMiddleware)
	e.GET("/backup/list", controller.ListBackups, authMiddleware)
}

// CreateBackup manually triggers a backup
func (controller *BackupController) CreateBackup(c echo.Context) error {
	err := controller.backupService.CreateBackup()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create backup: "+err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Backup created successfully",
	})
}

// ListBackups returns all available backups
func (controller *BackupController) ListBackups(c echo.Context) error {
	backups, err := controller.backupService.GetBackupsList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to list backups: "+err.Error()))
	}

	return c.JSON(http.StatusOK, backups)
}
