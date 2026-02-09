package dto

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