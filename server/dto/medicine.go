package dto

import "github.com/google/uuid"

type Schedule  struct {
	Time       string `json:"time" binding:"required"`
	DaysOfWeek []int  `json:"days_of_week" binding:"required"`
}

type CreateMedicineRequest struct {
	Name      string `json:"name" binding:"required"`
	Dosage    string `json:"dosage" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date"`
	Schedule  struct {
		Time       string `json:"time" binding:"required"`
		DaysOfWeek []int  `json:"days_of_week" binding:"required"`
	} `json:"schedule" binding:"required"`
}

type MedicineDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Dosage    string    `json:"dosage"`
	StartDate string    `json:"start_date"`
	EndDate   *string    `json:"end_date"`

	Schedule struct {
		Time       string `json:"time"`
		DaysOfWeek []int  `json:"days_of_week"`
	} `json:"schedule"`
}
