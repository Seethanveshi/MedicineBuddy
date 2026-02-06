package handler

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/mapper"
	"MedicineBuddy/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DoseHandler struct {
	doseService *service.DoseService
}

func NewDoseHandler(doseService *service.DoseService) *DoseHandler {
	return &DoseHandler{
		doseService: doseService,
	}
}

func (h *DoseHandler) TakeDose(c *gin.Context) {
	doseID, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dose id"})
		return
	}

	if err := h.doseService.MarkDoseTaken(c.Request.Context(), doseID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status": "Accepted"})
}

func (h *DoseHandler) SkipDose(c *gin.Context) {
	doseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid dose id"})
		return
	}

	if err := h.doseService.MarkDoseSkipped(c.Request.Context(), doseID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status": "Accepted"})
}

func (h *DoseHandler) GetToday(c *gin.Context) {
	doses, err := h.doseService.GetTodayDoses(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.DoseLogResponse, 0, len(doses))
	for _, d := range doses {
		responses = append(responses, mapper.ToDoseResponse(d))
	}

	c.JSON(200, responses)
}

func (h *DoseHandler) GetUpcoming(c *gin.Context) {
	days := 7
	if q := c.Query("days"); q != "" {
		if d, err := strconv.Atoi(q); err == nil {
			days = d
		}
	}

	doses, err := h.doseService.GetUpcomingDoses(c.Request.Context(), days)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	responses := make([]dto.DoseLogResponse, 0, len(doses))
	for _, d := range doses {
		responses = append(responses, mapper.ToDoseResponse(d))
	}

	c.JSON(200, responses)
}

func (h *DoseHandler) GetHistory(c *gin.Context) {
	limit := 50
	if q := c.Query("limit"); q != "" {
		if l, err := strconv.Atoi(q); err == nil {
			limit = l
		}
	}

	doses, err := h.doseService.GetDoseHistory(c.Request.Context(), limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	responses := make([]dto.DoseLogResponse, 0, len(doses))
	for _, d := range doses {
		responses = append(responses, mapper.ToDoseResponse(d))
	}

	c.JSON(200, responses)
}
