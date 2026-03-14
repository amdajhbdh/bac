package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	UserID     string                 `json:"user_id,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Status     string                 `json:"status"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type AuditLogger struct {
	pool    *pgxpool.Pool
	channel chan AuditEvent
}

func NewAuditLogger(pool *pgxpool.Pool) *AuditLogger {
	logger := &AuditLogger{
		pool:    pool,
		channel: make(chan AuditEvent, 1000),
	}

	go logger.processEvents()
	return logger
}

func (l *AuditLogger) Log(ctx context.Context, event AuditEvent) {
	event.ID = uuid.New().String()
	event.Timestamp = time.Now()

	select {
	case l.channel <- event:
	default:
		slog.Warn("audit channel full, dropping event")
	}
}

func (l *AuditLogger) processEvents() {
	for event := range l.channel {
		if err := l.saveEvent(event); err != nil {
			slog.Error("failed to save audit event", "error", err)
		}
	}
}

func (l *AuditLogger) saveEvent(event AuditEvent) error {
	details, _ := json.Marshal(event.Details)
	metadata, _ := json.Marshal(event.Metadata)

	_, err := l.pool.Exec(context.Background(), `
		INSERT INTO audit_log (
			id, timestamp, user_id, action, resource, resource_id,
			ip_address, user_agent, status, details, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, event.ID, event.Timestamp, event.UserID, event.Action, event.Resource,
		event.ResourceID, event.IPAddress, event.UserAgent, event.Status,
		details, metadata)

	return err
}

func (l *AuditLogger) Query(ctx context.Context, filter AuditFilter) ([]AuditEvent, error) {
	query := "SELECT id, timestamp, user_id, action, resource, resource_id, ip_address, user_agent, status FROM audit_log WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filter.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, filter.UserID)
		argNum++
	}
	if filter.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argNum)
		args = append(args, filter.Action)
		argNum++
	}
	if filter.Resource != "" {
		query += fmt.Sprintf(" AND resource = $%d", argNum)
		args = append(args, filter.Resource)
		argNum++
	}
	if !filter.StartTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp >= $%d", argNum)
		args = append(args, filter.StartTime)
		argNum++
	}
	if !filter.EndTime.IsZero() {
		query += fmt.Sprintf(" AND timestamp <= $%d", argNum)
		args = append(args, filter.EndTime)
		argNum++
	}

	query += " ORDER BY timestamp DESC LIMIT $100"

	rows, err := l.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []AuditEvent
	for rows.Next() {
		var e AuditEvent
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.UserID, &e.Action,
			&e.Resource, &e.ResourceID, &e.IPAddress, &e.UserAgent, &e.Status); err != nil {
			continue
		}
		events = append(events, e)
	}

	return events, nil
}

type AuditFilter struct {
	UserID    string
	Action    string
	Resource  string
	StartTime time.Time
	EndTime   time.Time
}

var ActionTypes = map[string]string{
	"user.login":           "Authentication",
	"user.logout":          "Authentication",
	"user.register":        "Authentication",
	"user.password_change": "Account",
	"question.create":      "Content",
	"question.update":      "Content",
	"question.delete":      "Content",
	"api.request":          "API",
	"admin.action":         "Admin",
}

func LogAction(ctx context.Context, logger *AuditLogger, action, resource, userID, ip string, status string) {
	logger.Log(ctx, AuditEvent{
		Action:    action,
		Resource:  resource,
		UserID:    userID,
		IPAddress: ip,
		Status:    status,
	})
}

func CreateAuditTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS audit_log (
			id UUID PRIMARY KEY,
			timestamp TIMESTAMPTZ NOT NULL,
			user_id VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			resource VARCHAR(100) NOT NULL,
			resource_id VARCHAR(255),
			ip_address INET,
			user_agent TEXT,
			status VARCHAR(50),
			details JSONB,
			metadata JSONB
		);

		CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_log(timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_log(user_id);
		CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);
		CREATE INDEX IF NOT EXISTS idx_audit_resource ON audit_log(resource);
	`)

	if err != nil {
		slog.Error("failed to create audit table", "error", err)
	}
	return err
}
