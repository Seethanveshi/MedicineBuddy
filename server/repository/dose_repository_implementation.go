package repository

import (
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

func (r *DoseRepositoryImple) Exists(ctx context.Context, medicineID uuid.UUID, scheduleAt time.Time) (bool, error) {
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

func (r *DoseRepositoryImple) Create(ctx context.Context, dose *model.DoseLog) error {
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

	_, err := r.db.ExecContext(
		ctx,
		query,
		model.DoseMissed,
		model.DosePending,
		now,
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
) ([]model.DoseLog, error) {

	query := `
		SELECT id, medicine_id, scheduled_at, status, taken_at
		FROM dose_logs
		WHERE scheduled_at >= $1 AND scheduled_at < $2
		ORDER BY scheduled_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doses []model.DoseLog

	for rows.Next() {
		var d model.DoseLog
		var takenAt sql.NullTime

		err := rows.Scan(
			&d.ID,
			&d.MedicineID,
			&d.ScheduleAt,
			&d.Status,
			&takenAt,
		)
		if err != nil {
			return nil, err
		}

		if takenAt.Valid {
			d.TakenAt = &takenAt.Time
		} else {
			d.TakenAt = nil
		}

		doses = append(doses, d)
	}

	return doses, nil
}

func (r *DoseRepositoryImple) GetDoseHistory(
	ctx context.Context,
	limit int,
) ([]model.DoseLog, error) {

	query := `
		SELECT id, medicine_id, scheduled_at, status, taken_at
		FROM dose_logs
		WHERE status IN ('taken','missed','skipped')
		ORDER BY scheduled_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doses []model.DoseLog

	for rows.Next() {
		var d model.DoseLog
		if err := rows.Scan(
			&d.ID,
			&d.MedicineID,
			&d.ScheduleAt,
			&d.Status,
			&d.TakenAt,
		); err != nil {
			return nil, err
		}
		doses = append(doses, d)
	}

	return doses, nil
}
