package handler

import (
	"MedicineBuddy/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(rs *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: rs}
}

func (h *ReportHandler) WeeklyDetailed(c *gin.Context) {

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	data, err := h.reportService.WeeklyDetailed(
		c.Request.Context(),
		patientID,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
