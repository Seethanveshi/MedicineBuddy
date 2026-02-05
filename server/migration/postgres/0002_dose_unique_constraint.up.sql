ALTER TABLE dose_logs
ADD CONSTRAINT unique_medicine_schedule
UNIQUE (medicine_id,scheduled_at);