package handler

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediTakerHandler struct {
	meditakerService *service.MediTakerService
}

func NewMediTakerHandler(mts *service.MediTakerService) *MediTakerHandler {
	return &MediTakerHandler{meditakerService: mts}
}

func (h *MediTakerHandler) CreateMediTaker(c *gin.Context) {
	var req dto.MediTakerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// until auth → use hardcoded or middleware
	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	err := h.meditakerService.CreateMediTaker(c.Request.Context(), patientID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "meditaker added"})
}

func (h *MediTakerHandler) List(c *gin.Context) {

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	data, err := h.meditakerService.ListByPatient(
		c.Request.Context(),
		patientID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *MediTakerHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	err = h.meditakerService.Delete(c.Request.Context(), id, patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *MediTakerHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var medireq dto.MediTakerRequest
	if err := c.ShouldBindJSON(&medireq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	err = h.meditakerService.Update(c.Request.Context(), id, patientID, medireq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Updated")
}
