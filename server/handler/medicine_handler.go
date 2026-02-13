package handler

import (
	"net/http"
	"time"

	"MedicineBuddy/dto"
	"MedicineBuddy/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MedicineHandler struct {
	medicineService *service.MedicineService
}

func NewMedicineHandler(ms *service.MedicineService) *MedicineHandler {
	return &MedicineHandler{medicineService: ms}
}

func (h *MedicineHandler) CreateMedicine(c *gin.Context) {
	var req dto.CreateMedicineRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, d := range req.Schedule.DaysOfWeek {
		if d < 0 || d > 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "days_of_week must be between 0 (Sunday) and 6 (Saturday)",
			})
			return
		}
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsedEnd, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil || parsedEnd.Before(startDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		endDate = &parsedEnd
	}

	timeOfDay, err := time.Parse("15:04", req.Schedule.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule time"})
		return
	}

	if err := h.medicineService.CreateMedicine(
		c.Request.Context(),
		req.Name,
		req.Dosage,
		startDate,
		endDate,
		timeOfDay,
		req.Schedule.DaysOfWeek,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *MedicineHandler) GetByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	data, err := h.medicineService.GetByID(
		c.Request.Context(),
		id,
		patientID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
