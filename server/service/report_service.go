package service

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/repository"
	"context"
	"time"

	"github.com/google/uuid"
)


type ReportService struct {
	doseRepository repository.DoseRepository
}

func NewReportService(repo repository.DoseRepository) *ReportService {
	return &ReportService{doseRepository: repo}
}


func (s *ReportService) WeeklyDetailed(
	ctx context.Context,
	patientID uuid.UUID,
) (dto.WeeklyDetailedReport, error) {

	now := time.Now()

	offset := (int(now.Weekday()) + 6) % 7
	start := time.Date(
		now.Year(), now.Month(), now.Day()-offset,
		0, 0, 0, 0, now.Location(),
	)
	end := start.AddDate(0, 0, 7)

	name, rows, err := s.doseRepository.WeeklyDetailed(ctx, patientID, start, end)
	if err != nil {
		return dto.WeeklyDetailedReport{}, err
	}

	result := dto.WeeklyDetailedReport{
		PatientName: name,
	}

	dayMap := map[string][]dto.DoseDetail{}

	for _, r := range rows {
		date := r.ScheduledAt.Format("2006-01-02")
		dayMap[date] = append(dayMap[date], r)
	}

	for date, doses := range dayMap {
		result.Days = append(result.Days, dto.DayReport{
			Date:  date,
			Doses: doses,
		})
	}

	return result, nil
}

