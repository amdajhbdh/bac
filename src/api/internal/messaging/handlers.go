package messaging

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebhookPayload represents incoming webhook data
type WebhookPayload struct {
	Message string                 `json:"message"`
	To      string                 `json:"to,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Handler holds dependencies for messaging handlers
type Handler struct {
	// Future: inject Telegram/WhatsApp clients
}

// NewHandler creates a new messaging handler
func NewHandler() *Handler {
	return &Handler{}
}

// TelegramWebhook handles incoming Telegram updates
func (h *Handler) TelegramWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Telegram] Received message: %s", payload.Message)
	// TODO: forward to Telegram Bot API

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// WhatsAppWebhook handles incoming WhatsApp messages
func (h *Handler) WhatsAppWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[WhatsApp] Received message: %s", payload.Message)
	// TODO: forward to WhatsApp Business API

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// Health checks handler status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "messaging",
		"status":  "ok",
	})
}