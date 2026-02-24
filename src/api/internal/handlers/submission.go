package handlers

import (
	"github.com/bac-unified/api/internal/services"

	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	submissionService *services.SubmissionService
}

func NewSubmissionHandler(submissionService *services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService}
}

func (h *SubmissionHandler) SubmitImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "image required"})
		return
	}
	defer file.Close()

	data := make([]byte, header.Size)
	file.Read(data)

	subject := c.PostForm("subject")
	resp, err := h.submissionService.SubmitImage(data, header.Filename, subject)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(202, resp)
}

func (h *SubmissionHandler) SubmitPDF(c *gin.Context) {
	file, header, err := c.Request.FormFile("pdf")
	if err != nil {
		c.JSON(400, gin.H{"error": "pdf required"})
		return
	}
	defer file.Close()

	data := make([]byte, header.Size)
	file.Read(data)

	subject := c.PostForm("subject")
	resp, err := h.submissionService.SubmitPDF(data, header.Filename, subject)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(202, resp)
}

func (h *SubmissionHandler) SubmitURL(c *gin.Context) {
	var req struct {
		URL     string `json:"url"`
		Subject string `json:"subject"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.submissionService.SubmitURL(req.URL, req.Subject)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(202, resp)
}

func (h *SubmissionHandler) GetStatus(c *gin.Context) {
	jobID := c.Param("job_id")
	resp, err := h.submissionService.GetStatus(jobID)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, resp)
}
