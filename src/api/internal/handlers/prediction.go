package handlers

import (
	"github.com/bac-unified/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PredictionHandler struct {
	predictionService *services.PredictionService
}

func NewPredictionHandler(predictionService *services.PredictionService) *PredictionHandler {
	return &PredictionHandler{predictionService: predictionService}
}

func (h *PredictionHandler) ListPredictions(c *gin.Context) {
	_ = c.Query("subject_id")
	_ = c.Query("exam_year")

	// Simplified - would query database
	c.JSON(200, gin.H{"predictions": []interface{}{}})
}

func (h *PredictionHandler) GetPrediction(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	prediction, err := h.predictionService.GetPrediction(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, prediction)
}

func (h *PredictionHandler) GetBySubject(c *gin.Context) {
	subjectID := 1
	examYear := 2027

	predictions, err := h.predictionService.GetBySubject(subjectID, examYear)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, predictions)
}

func (h *PredictionHandler) GetLatest(c *gin.Context) {
	predictions, err := h.predictionService.GetLatest()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, predictions)
}
