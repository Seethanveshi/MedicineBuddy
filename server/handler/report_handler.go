package handler

import (
	"MedicineBuddy/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportService    *service.ReportService
	mediTakerService *service.MediTakerService
	emailService     *service.SMTPEmailService
}

func NewReportHandler(rs *service.ReportService, ms *service.MediTakerService, es *service.SMTPEmailService) *ReportHandler {
	return &ReportHandler{
		reportService:    rs,
		mediTakerService: ms,
		emailService:     es,
	}
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

func (h *ReportHandler) SendTest(c *gin.Context) {

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	report, err := h.reportService.WeeklyDetailed(
		c.Request.Context(),
		patientID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meditakers, err := h.mediTakerService.ListByPatient(
		c.Request.Context(),
		patientID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	body := h.emailService.BuildWeeklyEmail(report)

	for _, m := range meditakers {
		err = h.emailService.Send(
			m.Email,
			"Weekly medicine report",
			body,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "emails sent"})
}
