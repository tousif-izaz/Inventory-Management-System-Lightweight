package backup

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
)

type BackupConfig struct {
	SourcePath     string        // Path to the SQLite database file
	BackupDir      string        // Directory where backups will be stored
	RetentionDays  int           // Number of days to keep backups
	BackupInterval time.Duration // Interval between automatic backups
}

type BackupService struct {
	config BackupConfig
	ticker *time.Ticker
	done   chan bool
}

// NewBackupService creates a new backup service
func NewBackupService(config BackupConfig) *BackupService {
	return &BackupService{
		config: config,
		done:   make(chan bool),
	}
}

// Start begins the scheduled backup process
func (s *BackupService) Start() {
	// Ensure backup directory exists
	if err := os.MkdirAll(s.config.BackupDir, 0755); err != nil {
		log.Errorf("Failed to create backup directory: %v", err)
		return
	}

	// Create initial backup on startup
	log.Info("Creating initial backup on startup...")
	if err := s.CreateBackup(); err != nil {
		log.Errorf("Initial backup failed: %v", err)
	} else {
		log.Info("Initial backup created successfully")
	}

	// Start scheduled backups
	s.ticker = time.NewTicker(s.config.BackupInterval)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Info("Running scheduled backup...")
				if err := s.CreateBackup(); err != nil {
					log.Errorf("Scheduled backup failed: %v", err)
				} else {
					log.Info("Scheduled backup created successfully")
				}

				// Clean up old backups
				if err := s.CleanOldBackups(); err != nil {
					log.Errorf("Backup cleanup failed: %v", err)
				}
			case <-s.done:
				return
			}
		}
	}()

	log.Infof("Backup service started. Backups will run every %v", s.config.BackupInterval)
}

// Stop stops the scheduled backup process
func (s *BackupService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true

	// Create final backup on shutdown
	log.Info("Creating final backup on shutdown...")
	if err := s.CreateBackup(); err != nil {
		log.Errorf("Final backup failed: %v", err)
	} else {
		log.Info("Final backup created successfully")
	}
}

// CreateBackup creates a new backup of the database
func (s *BackupService) CreateBackup() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFilename := fmt.Sprintf("inventory_backup_%s.db", timestamp)
	backupPath := filepath.Join(s.config.BackupDir, backupFilename)

	// Open source database file
	src, err := os.Open(s.config.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer src.Close()

	// Create backup file
	dst, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	// Copy database file
	bytesWritten, err := io.Copy(dst, src)
	if err != nil {
		// Remove incomplete backup file
		os.Remove(backupPath)
		return fmt.Errorf("failed to copy database: %w", err)
	}

	log.Infof("Backup created: %s (%d bytes)", backupFilename, bytesWritten)
	return nil
}

// CleanOldBackups removes backups older than retention period
func (s *BackupService) CleanOldBackups() error {
	cutoffTime := time.Now().AddDate(0, 0, -s.config.RetentionDays)
	deletedCount := 0

	err := filepath.Walk(s.config.BackupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only consider backup files
		if filepath.Ext(path) == ".db" && info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				log.Errorf("Failed to remove old backup %s: %v", path, err)
			} else {
				deletedCount++
				log.Infof("Removed old backup: %s", filepath.Base(path))
			}
		}

		return nil
	})

	if deletedCount > 0 {
		log.Infof("Cleaned up %d old backup(s)", deletedCount)
	}

	return err
}

// GetBackupsList returns a list of all available backups
func (s *BackupService) GetBackupsList() ([]BackupInfo, error) {
	var backups []BackupInfo

	err := filepath.Walk(s.config.BackupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".db" {
			backups = append(backups, BackupInfo{
				Filename:  filepath.Base(path),
				Path:      path,
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
		}

		return nil
	})

	return backups, err
}

// RestoreBackup restores a backup file to the main database
// WARNING: This will overwrite the current database
func (s *BackupService) RestoreBackup(backupFilename string, db *sql.DB) error {
	backupPath := filepath.Join(s.config.BackupDir, backupFilename)

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFilename)
	}

	// Close all database connections before restore
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	// Backup current database before restore (safety measure)
	safetyBackupPath := filepath.Join(s.config.BackupDir, fmt.Sprintf("pre_restore_%s.db", time.Now().Format("2006-01-02_15-04-05")))
	if err := copyFile(s.config.SourcePath, safetyBackupPath); err != nil {
		log.Warnf("Failed to create safety backup: %v", err)
	}

	// Restore the backup
	if err := copyFile(backupPath, s.config.SourcePath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	log.Infof("Database restored from backup: %s", backupFilename)
	return nil
}

// BackupInfo contains information about a backup file
type BackupInfo struct {
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// Helper function to copy files
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
