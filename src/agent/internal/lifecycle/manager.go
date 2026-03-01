package lifecycle

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type LifecyclePolicy struct {
	Bucket         string
	Prefix         string
	TransitionDays int
	ArchiveDays    int
	DeleteDays     int
}

type LifecycleManager struct {
	policies []LifecyclePolicy
}

func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		policies: []LifecyclePolicy{},
	}
}

func (m *LifecycleManager) AddPolicy(policy LifecyclePolicy) {
	m.policies = append(m.policies, policy)
	slog.Info("added lifecycle policy", "prefix", policy.Prefix)
}

func (m *LifecycleManager) ApplyToS3(ctx context.Context) error {
	slog.Info("applying lifecycle policies to S3", "count", len(m.policies))

	for _, policy := range m.policies {
		if err := m.applyPolicy(ctx, policy); err != nil {
			slog.Warn("failed to apply policy", "prefix", policy.Prefix, "error", err)
		}
	}

	return nil
}

func (m *LifecycleManager) applyPolicy(ctx context.Context, policy LifecyclePolicy) error {
	_ = fmt.Sprintf("lifecycle-%s", policy.Prefix)

	slog.Info("would apply lifecycle rule",
		"bucket", policy.Bucket,
		"prefix", policy.Prefix,
		"transition_days", policy.TransitionDays,
		"archive_days", policy.ArchiveDays,
		"delete_days", policy.DeleteDays)

	return nil
}

func GetDefaultPolicies() []LifecyclePolicy {
	return []LifecyclePolicy{
		{
			Bucket:         os.Getenv("GARAGE_BUCKET"),
			Prefix:         "uploads/",
			TransitionDays: 30,
			ArchiveDays:    90,
			DeleteDays:     365,
		},
		{
			Bucket:         os.Getenv("GARAGE_BUCKET"),
			Prefix:         "animations/",
			TransitionDays: 7,
			ArchiveDays:    30,
			DeleteDays:     180,
		},
		{
			Bucket:         os.Getenv("GARAGE_BUCKET"),
			Prefix:         "cache/",
			TransitionDays: 1,
			ArchiveDays:    7,
			DeleteDays:     30,
		},
		{
			Bucket:         os.Getenv("GARAGE_BUCKET"),
			Prefix:         "backups/",
			TransitionDays: 90,
			ArchiveDays:    365,
			DeleteDays:     2555,
		},
	}
}

func InitDefaultLifecycle() {
	manager := NewLifecycleManager()
	for _, policy := range GetDefaultPolicies() {
		manager.AddPolicy(policy)
	}
	slog.Info("lifecycle policies initialized")
}

type CleanupResult struct {
	FilesDeleted int
	SpaceFreed   int64
}

func (m *LifecycleManager) Cleanup(ctx context.Context, bucket, prefix string, olderThan time.Duration) (*CleanupResult, error) {
	slog.Info("running cleanup", "bucket", bucket, "prefix", prefix, "older_than", olderThan)

	result := &CleanupResult{
		FilesDeleted: 0,
		SpaceFreed:   0,
	}

	return result, nil
}
