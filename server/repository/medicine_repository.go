package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/google/uuid"
)

type MedicineRepository struct {
	db *sql.DB
}

func NewMedicineRepository(db *sql.DB) *MedicineRepository {
	return &MedicineRepository{db: db}
}

func (r *MedicineRepository) CreateMedicineWithSchedule(ctx context.Context, medicine *model.Medicine, schedule *model.Schedule) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		INSERT INTO medicines(id, user_id, name, dosage, start_date, end_date) 
		VALUES($1, $2, $3, $4, $5, $6)	
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		medicine.ID,
		medicine.UserID,
		medicine.Name,
		medicine.Dosage,
		medicine.StartDate,
		medicine.EndDate,
	)
	if err != nil {
		return err
	}

	query1 := `
		INSERT INTO schedules(id, medicine_id, time_of_day, days_of_week)
		VALUES($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(
		ctx,
		query1,
		schedule.ID,
		schedule.MedicineID,
		schedule.TimeOfDay,
		schedule.DaysOfWeek,
	)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *MedicineRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
	patientID uuid.UUID,
) (dto.MedicineDetailResponse, error) {

	query := `
	SELECT 
		m.id,
		m.name,
		m.dosage,
		m.start_date,
		m.end_date,
		s.time_of_day,
		s.days_of_week
	FROM medicines m
	JOIN schedules s ON s.medicine_id = m.id
	WHERE m.id = $1
	AND m.user_id = $2
	`

	var res dto.MedicineDetailResponse
	var days pq.Int64Array

	err := r.db.QueryRowContext(ctx, query, id, patientID).
		Scan(
			&res.ID,
			&res.Name,
			&res.Dosage,
			&res.StartDate,
			&res.EndDate,
			&res.Schedule.Time,
			&days,
		)

	if err != nil {
		return res, err
	}

	res.Schedule.DaysOfWeek = make([]int, len(days))
	for i, d := range days {
		res.Schedule.DaysOfWeek[i] = int(d)
	}

	return res, nil
}
