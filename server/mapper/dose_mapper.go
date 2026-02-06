package mapper

import (
	"time"

	"MedicineBuddy/dto"
	"MedicineBuddy/model"
)

func ToDoseResponse(d model.DoseLog) dto.DoseLogResponse {

	var takenAt *string
	if d.TakenAt != nil {
		t := d.TakenAt.UTC().Format(time.RFC3339)
		takenAt = &t
	}

	return dto.DoseLogResponse{
		ID:          d.ID.String(),
		MedicineID:  d.MedicineID.String(),
		ScheduledAt: d.ScheduleAt.UTC().Format(time.RFC3339),
		Status:      d.Status,
		TakenAt:     takenAt,
	}
}
