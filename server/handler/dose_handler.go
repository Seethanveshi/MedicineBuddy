package handler

import (
	"MedicineBuddy/service"
	"net/http"

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
