package handlers

import (
	"github.com/bac-unified/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SolverHandler struct {
	solverService *services.SolverService
}

func NewSolverHandler(solverService *services.SolverService) *SolverHandler {
	return &SolverHandler{solverService: solverService}
}

func (h *SolverHandler) Solve(c *gin.Context) {
	var req services.SolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.solverService.Solve(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *SolverHandler) GetSteps(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	steps, err := h.solverService.GetSteps(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, steps)
}

func (h *SolverHandler) GetAnimation(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	animation, err := h.solverService.GetAnimation(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, gin.H{"animation_url": animation})
}
