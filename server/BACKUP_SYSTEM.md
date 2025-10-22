# Database Backup System

## Overview
Automated backup system for the SQLite database with scheduled backups, retention policies, and manual backup triggers.

## Features

### 1. **Automatic Scheduled Backups**
- Backups run automatically every 24 hours
- First backup created immediately on server startup
- Final backup created on graceful server shutdown

### 2. **Backup Retention**
- Keeps backups for 7 days by default
- Old backups are automatically cleaned up
- Configurable retention period

### 3. **Manual Backup Trigger**
- API endpoint to create backups on-demand
- Useful before major operations or updates

### 4. **Backup Information**
- List all available backups with metadata
- View backup size and creation time

## Configuration

Located in `server/cmd/imsapi/main.go`:

```go
backupConfig := backup.BackupConfig{
    SourcePath:     "data/inventory.db",  // Database file path
    BackupDir:      "data/backups",        // Backup storage directory
    RetentionDays:  7,                     // Keep backups for 7 days
    BackupInterval: 24 * time.Hour,        // Backup every 24 hours
}
```

### Configurable Parameters:
- **SourcePath**: Path to the SQLite database file
- **BackupDir**: Directory where backups will be stored
- **RetentionDays**: Number of days to keep old backups (default: 7)
- **BackupInterval**: Time between automatic backups (default: 24 hours)

## API Endpoints

### Create Manual Backup
```
POST /backup
Authorization: Required (JWT token)
```

**Response:**
```json
{
  "message": "Backup created successfully"
}
```

### List Available Backups
```
GET /backup/list
Authorization: Required (JWT token)
```

**Response:**
```json
[
  {
    "filename": "inventory_backup_2024-01-15_14-30-00.db",
    "path": "data/backups/inventory_backup_2024-01-15_14-30-00.db",
    "size": 8192,
    "created_at": "2024-01-15T14:30:00Z"
  }
]
```

## Backup File Naming Convention

Backups are named with timestamps for easy identification:
```
inventory_backup_YYYY-MM-DD_HH-MM-SS.db
```

Example: `inventory_backup_2024-01-15_14-30-00.db`

## Backup Schedule

1. **On Server Start**: Immediate backup created
2. **During Runtime**: Backup every 24 hours
3. **On Server Stop**: Final backup before shutdown
4. **Manual**: Anytime via API endpoint

## Storage Location

Backups are stored in: `data/backups/`

## Backup Process

1. Database file is safely copied while server is running
2. Backup is saved with timestamp
3. Old backups (older than retention period) are automatically deleted
4. Logs show backup creation and cleanup activities

## Logs

The backup system logs the following events:
- Backup service started
- Backup created successfully
- Old backups cleaned up
- Any errors during backup process

Example logs:
```
Backup service started. Backups will run every 24h0m0s
Creating initial backup on startup...
Initial backup created successfully
Backup created: inventory_backup_2024-01-15_14-30-00.db (8192 bytes)
```

## Restoration Process

To restore from a backup:

1. **Stop the server**
2. **Backup current database** (safety measure)
3. **Copy backup file** to replace current database:
   ```bash
   copy "data\backups\inventory_backup_2024-01-15_14-30-00.db" "data\inventory.db"
   ```
4. **Restart the server**

**Note**: Manual restoration is currently the recommended approach. Automatic restoration via API can be added if needed.

## Best Practices

1. **Regular Monitoring**: Check backup logs to ensure backups are running
2. **Test Restores**: Periodically test backup restoration in a test environment
3. **External Backups**: Consider copying backups to external storage or cloud
4. **Before Updates**: Create manual backup before major system updates
5. **Disk Space**: Monitor backup directory to ensure sufficient disk space

## Troubleshooting

### Backup Failed
- Check write permissions for backup directory
- Ensure sufficient disk space
- Review server logs for detailed error messages

### Backups Not Running
- Verify backup service started (check logs)
- Check backup interval configuration
- Ensure server is running continuously

### Old Backups Not Deleted
- Check retention policy configuration
- Verify write permissions for backup directory
- Review cleanup logs

## Security

- Backup endpoints require authentication
- Only authenticated users can create or list backups
- Backup files are stored locally with restricted access
- Consider encrypting backups for sensitive data

## Future Enhancements

Potential improvements:
- Cloud backup integration (AWS S3, Google Drive)
- Automated restore via API
- Backup compression to save space
- Email notifications on backup failure
- Differential/incremental backups
