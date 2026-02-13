package repository

import (
	"MedicineBuddy/dto"
	"MedicineBuddy/model"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type DoseRepositoryImple struct {
	db *sql.DB
}

func NewDoseRepository(db *sql.DB) DoseRepository {
	return &DoseRepositoryImple{db: db}
}

func (r *DoseRepositoryImple) Exists(ctx context.Context, tx *sql.Tx, medicineID uuid.UUID, scheduleAt time.Time) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM dose_logs
			WHERE medicine_id = $1 AND scheduled_at = $2
		)
	`

	err := r.db.QueryRowContext(ctx, query, medicineID, scheduleAt).Scan(&exists)
	return exists, err
}

func (r *DoseRepositoryImple) Create(ctx context.Context, tx *sql.Tx, dose *model.DoseLog) error {
	query := `
		INSERT INTO dose_logs (id, medicine_id, scheduled_at, status)
		VALUES($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		dose.ID,
		dose.MedicineID,
		dose.ScheduleAt,
		dose.Status,
	)

	return err
}

func (r *DoseRepositoryImple) MarkMissedDose(ctx context.Context, now time.Time) error {
	query := `
		UPDATE dose_logs
		SET status = $1
		WHERE status = $2 AND scheduled_at < $3
	`
	delay := 10 * time.Minute
	_, err := r.db.ExecContext(
		ctx,
		query,
		model.DoseMissed,
		model.DosePending,
		now.Add(-delay),
	)

	return err
}

func (r *DoseRepositoryImple) UpdateDoseStatus(ctx context.Context, doseID uuid.UUID, fromStatus string, toStatus string, takenAt *time.Time) error {
	query := `
		UPDATE dose_logs
		SET status = $1, taken_at = $2
		WHERE id = $3 AND status = $4
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		toStatus,
		takenAt,
		doseID,
		fromStatus,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("invalid dose state transition")
	}

	return nil
}

func (r *DoseRepositoryImple) GetDosesBetween(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]dto.DoseLogResponse, error) {

	query := `
		SELECT dl.id, m.name, m.dosage, dl.medicine_id, dl.scheduled_at, dl.status, dl.taken_at
		FROM dose_logs dl
		LEFT JOIN medicines m ON dl.medicine_id = m.id
		WHERE dl.scheduled_at >= $1 AND dl.scheduled_at < $2
		ORDER BY dl.scheduled_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doses []dto.DoseLogResponse

	for rows.Next() {
		var d dto.DoseLogResponse
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Dosage,
			&d.MedicineID,
			&d.ScheduledAt,
			&d.Status,
			&d.TakenAt,
		)

		if err != nil {
			return nil, err
		}

		doses = append(doses, d)
	}

	return doses, nil
}

func (r *DoseRepositoryImple) GetDoseHistory(
	ctx context.Context,
	limit int,
) ([]dto.DoseLogResponse, error) {

	query := `
		SELECT dl.id, m.name, m.dosage, dl.medicine_id, dl.scheduled_at, dl.status, dl.taken_at
		FROM dose_logs dl
		LEFT JOIN medicines m ON m.id = dl.medicine_id
		WHERE dl.status IN ('taken','missed','skipped')
		ORDER BY dl.scheduled_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doses []dto.DoseLogResponse

	for rows.Next() {
		var d dto.DoseLogResponse
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Dosage,
			&d.MedicineID,
			&d.ScheduledAt,
			&d.Status,
			&d.TakenAt,
		); err != nil {
			return nil, err
		}
		doses = append(doses, d)
	}

	return doses, nil
}

func (r *DoseRepositoryImple) GetDosesByDate(ctx context.Context, date time.Time) ([]dto.DoseLogResponse, error) {

	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	query := `
		SELECT dl.id, m.name, m.dosage,dl.medicine_id, dl.scheduled_at, dl.status, dl.taken_at
		FROM dose_logs dl
		LEFT JOIN medicines m ON m.id = dl.medicine_id 
		WHERE dl.scheduled_at >= $1
		AND dl.scheduled_at < $2
		ORDER BY dl.scheduled_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doses []dto.DoseLogResponse

	for rows.Next() {
		var d dto.DoseLogResponse

		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Dosage,
			&d.MedicineID,
			&d.ScheduledAt,
			&d.Status,
			&d.TakenAt,
		)
		if err != nil {
			return nil, err
		}

		doses = append(doses, d)
	}

	return doses, nil
}

func (r *DoseRepositoryImple) WeeklyDetailed(
	ctx context.Context,
	patientID uuid.UUID,
	start time.Time,
	end time.Time,
) (string, []dto.DoseDetail, error) {

	query := `
	SELECT 
		u.user_name,
		m.name,
		m.dosage,
		dl.scheduled_at,
		dl.status,
		dl.taken_at
	FROM dose_logs dl
	JOIN medicines m ON m.id = dl.medicine_id
	JOIN users u ON u.id = m.user_id
	WHERE m.user_id = $1
	AND dl.scheduled_at >= $2
	AND dl.scheduled_at < $3
	ORDER BY dl.scheduled_at
	`

	rows, err := r.db.QueryContext(ctx, query, patientID, start, end)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var (
		patientName string
		list        []dto.DoseDetail
	)

	for rows.Next() {
		var d dto.DoseDetail

		err := rows.Scan(
			&patientName,
			&d.MedicineName,
			&d.Dosage,
			&d.ScheduledAt,
			&d.Status,
			&d.TakenAt,
		)
		if err != nil {
			return "", nil, err
		}

		list = append(list, d)
	}

	return patientName, list, nil
}
