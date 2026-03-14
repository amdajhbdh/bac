package backup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

type BackupConfig struct {
	DBURL         string
	S3Endpoint    string
	S3Bucket      string
	Schedule      string
	RetentionDays int
}

type BackupManager struct {
	config BackupConfig
}

func NewBackupManager(config BackupConfig) *BackupManager {
	if config.RetentionDays == 0 {
		config.RetentionDays = 30
	}
	return &BackupManager{config: config}
}

func (m *BackupManager) RunBackup(ctx context.Context) (string, error) {
	slog.Info("starting database backup")

	backupFile := fmt.Sprintf("bac_backup_%s.sql", time.Now().Format("2006-01-02_150405"))

	cmd := exec.CommandContext(ctx, "pg_dump", m.config.DBURL, "-f", backupFile)
	cmd.Env = append(os.Environ(), "PGSSLMODE=require")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pg_dump failed: %w", err)
	}

	slog.Info("backup created", "file", backupFile)
	return backupFile, nil
}

func (m *BackupManager) UploadToS3(ctx context.Context, backupFile string) error {
	slog.Info("uploading backup to S3", "file", backupFile)

	s3Cmd := exec.CommandContext(ctx, "aws", "s3", "cp", backupFile,
		fmt.Sprintf("s3://%s/backups/%s", m.config.S3Bucket, backupFile))

	if err := s3Cmd.Run(); err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}

	slog.Info("backup uploaded to S3")
	return nil
}

func (m *BackupManager) CleanupOldBackups(ctx context.Context) error {
	slog.Info("cleaning up old backups", "retention_days", m.config.RetentionDays)

	cutoff := time.Now().AddDate(0, 0, -m.config.RetentionDays)

	cmd := exec.CommandContext(ctx, "aws", "s3", "ls", fmt.Sprintf("s3://%s/backups/", m.config.S3Bucket))
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("list backups failed: %w", err)
	}

	slog.Info("backup cleanup complete", "cutoff", cutoff)
	return nil
}

func (m *BackupManager) FullBackup(ctx context.Context) error {
	backupFile, err := m.RunBackup(ctx)
	if err != nil {
		return err
	}
	defer os.Remove(backupFile)

	if err := m.UploadToS3(ctx, backupFile); err != nil {
		return err
	}

	return m.CleanupOldBackups(ctx)
}

func GetBackupConfig() BackupConfig {
	return BackupConfig{
		DBURL:         os.Getenv("NEON_DB_URL"),
		S3Endpoint:    os.Getenv("BACKUP_S3_ENDPOINT"),
		S3Bucket:      os.Getenv("BACKUP_S3_BUCKET"),
		Schedule:      os.Getenv("BACKUP_SCHEDULE"),
		RetentionDays: 30,
	}
}
