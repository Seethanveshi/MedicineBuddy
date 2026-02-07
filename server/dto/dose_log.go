package dto

type DoseLogResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	MedicineID  string  `json:"medicine_id"`
	ScheduledAt string  `json:"scheduled_at"`
	Status      string  `json:"status"`
	TakenAt     *string `json:"taken_at"`
}
