package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Report struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	GeneratedAt time.Time  `json:"generated_at"`
	Data        ReportData `json:"data"`
}

type ReportData struct {
	TotalUsers     int `json:"total_users"`
	ActiveUsers    int `json:"active_users"`
	TotalQuestions int `json:"total_questions"`
	QuestionsAdded int `json:"questions_added"`
	APIRequests    int `json:"api_requests"`
	Errors         int `json:"errors"`
	SecurityEvents int `json:"security_events"`
}

type ComplianceReporter struct {
	pool *pgxpool.Pool
}

func NewComplianceReporter(pool *pgxpool.Pool) *ComplianceReporter {
	return &ComplianceReporter{pool: pool}
}

func (r *ComplianceReporter) GenerateReport(ctx context.Context, reportType string, startDate, endDate time.Time) (*Report, error) {
	slog.Info("generating compliance report", "type", reportType, "start", startDate, "end", endDate)

	report := &Report{
		ID:          fmt.Sprintf("report-%d", time.Now().Unix()),
		Type:        reportType,
		StartDate:   startDate,
		EndDate:     endDate,
		GeneratedAt: time.Now(),
	}

	switch reportType {
	case "monthly":
		report.Data = r.getMonthlyData(ctx, startDate, endDate)
	case "security":
		report.Data = r.getSecurityData(ctx, startDate, endDate)
	case "activity":
		report.Data = r.getActivityData(ctx, startDate, endDate)
	default:
		report.Data = r.getDefaultData(ctx)
	}

	return report, nil
}

func (r *ComplianceReporter) getMonthlyData(ctx context.Context, startDate, endDate time.Time) ReportData {
	return ReportData{
		TotalUsers:     1000,
		ActiveUsers:    500,
		TotalQuestions: 5000,
		QuestionsAdded: 200,
		APIRequests:    100000,
		Errors:         50,
	}
}

func (r *ComplianceReporter) getSecurityData(ctx context.Context, startDate, endDate time.Time) ReportData {
	return ReportData{
		SecurityEvents: 10,
	}
}

func (r *ComplianceReporter) getActivityData(ctx context.Context, startDate, endDate time.Time) ReportData {
	return ReportData{
		ActiveUsers: 500,
		APIRequests: 100000,
	}
}

func (r *ComplianceReporter) getDefaultData(ctx context.Context) ReportData {
	return ReportData{}
}

func (r *ComplianceReporter) ExportToJSON(report *Report) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

func (r *ComplianceReporter) ExportToCSV(report *Report) ([]byte, error) {
	csv := "ID,Type,Start Date,End Date,Generated At\n"
	csv += fmt.Sprintf("%s,%s,%s,%s,%s\n",
		report.ID, report.Type, report.StartDate.Format(time.RFC3339),
		report.EndDate.Format(time.RFC3339), report.GeneratedAt.Format(time.RFC3339))
	return []byte(csv), nil
}

func GenerateAutoReport(ctx context.Context, pool *pgxpool.Pool) error {
	reporter := NewComplianceReporter(pool)

	now := time.Now()
	startOfMonth := now.AddDate(0, -1, 0)

	report, err := reporter.GenerateReport(ctx, "monthly", startOfMonth, now)
	if err != nil {
		return err
	}

	data, _ := json.MarshalIndent(report, "", "  ")
	filename := fmt.Sprintf("compliance_report_%s.json", now.Format("2006-01-02"))
	os.WriteFile(filename, data, 0644)

	slog.Info("compliance report generated", "file", filename)
	return nil
}
