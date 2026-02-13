package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *MedicineRepository) DeleteFutureByMedicine(
	ctx context.Context,
	tx *sql.Tx,
	medicineID uuid.UUID,
	from time.Time,
) error {

	query := `
		DELETE FROM dose_logs
		WHERE medicine_id = $1
		AND scheduled_at > $2
		AND status = 'pending'
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		medicineID,
		from,
	)

	return err
}

func (r *MedicineRepository) UpdateTx(
	ctx context.Context,
	tx *sql.Tx,
	medicineID uuid.UUID,
	patientID uuid.UUID,
	req dto.UpdateReq,
) error {

	query := `
		UPDATE medicines
		SET name = $1,
		    dosage = $2,
		    start_date = $3,
		    end_date = $4
		WHERE id = $5
		AND user_id = $6
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		req.Name,
		req.Dosage,
		req.StartDate,
		req.EndDate,
		medicineID,
		patientID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("medicine not found")
	}

	return nil
}

func (r *MedicineRepository) UpdateScheduleByMedicineTx(
	ctx context.Context,
	tx *sql.Tx,
	medicineID uuid.UUID,
	time time.Time,
	dayOfWeek []int,
) error {

	query := `
		UPDATE schedules
		SET time_of_day = $1,
		    days_of_week = $2
		WHERE medicine_id = $3
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		time,
		pq.Array(dayOfWeek),
		medicineID,
	)

	return err
}

func (r *MedicineRepository) GetScheduleByMedicineTx(
	ctx context.Context,
	tx *sql.Tx,
	medicineID uuid.UUID,
) (*model.Schedule, error) {

	query := `
		SELECT id, medicine_id, time_of_day, days_of_week
		FROM schedules
		WHERE medicine_id = $1
	`

	var s model.Schedule
	var days pq.Int64Array

	err := tx.QueryRowContext(ctx, query, medicineID).
		Scan(
			&s.ID,
			&s.MedicineID,
			&s.TimeOfDay,
			&days,
		)
	if err != nil {
		return nil, err
	}

	s.DaysOfWeek = make([]int, len(days))
	for i, d := range days {
		s.DaysOfWeek[i] = int(d)
	}

	return &s, nil
}

func (r *MedicineRepository) GetModelByIDTx(
	ctx context.Context,
	tx *sql.Tx,
	medicineID uuid.UUID,
) (*model.Medicine, error) {

	query := `
		SELECT id, user_id, name, dosage, start_date, end_date
		FROM medicines
		WHERE id = $1
	`

	var m model.Medicine

	err := tx.QueryRowContext(ctx, query, medicineID).
		Scan(
			&m.ID,
			&m.UserID,
			&m.Name,
			&m.Dosage,
			&m.StartDate,
			&m.EndDate,
		)

	if err != nil {
		return nil, err
	}

	return &m, nil
}
