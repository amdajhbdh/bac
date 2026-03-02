package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bac-unified/agent/internal/ocr"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type Server struct {
	addr  string
	stats *Stats
}

type Stats struct {
	requestsTotal   atomic.Int64
	errorsTotal     atomic.Int64
	processingMsSum atomic.Int64
}

func NewServer(addr string) *Server {
	return &Server{
		addr:  addr,
		stats: &Stats{},
	}
}

func (s *Server) HandleOCR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ImageData string `json:"image_data"`
		Language  string `json:"language"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()
	ctx := context.Background()
	result := ocr.ProcessImage(ctx, req.ImageData)
	processingTime := time.Since(start).Milliseconds()

	if result.Error != "" {
		s.stats.errorsTotal.Add(1)
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error})
		return
	}

	s.stats.requestsTotal.Add(1)
	s.stats.processingMsSum.Add(processingTime)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"text":               result.Text,
		"confidence":         result.Confidence,
		"source":             result.Source,
		"processing_time_ms": processingTime,
	})
}

func (s *Server) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	requests := s.stats.requestsTotal.Load()
	errors := s.stats.errorsTotal.Load()
	sum := s.stats.processingMsSum.Load()

	avgMs := 0.0
	if requests > 0 {
		avgMs = float64(sum) / float64(requests)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"requests_total":    requests,
		"errors_total":      errors,
		"avg_processing_ms": avgMs,
	})
}

func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "healthy",
		"requests": s.stats.requestsTotal.Load(),
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	addr := os.Getenv("OCR_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8083"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := NewServer(addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/ocr", server.HandleOCR)
	mux.HandleFunc("/metrics", server.HandleMetrics)
	mux.HandleFunc("/health", server.HandleHealth)

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		logger.Info("OCR server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
