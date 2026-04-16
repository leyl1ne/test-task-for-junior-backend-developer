ALTER TABLE tasks
ADD COLUMN schedule_id BIGINT,
ADD COLUMN date DATE;

ALTER TABLE tasks
ADD CONSTRAINT fk_schedule
FOREIGN KEY (schedule_id) REFERENCES schedules(id)
ON DELETE CASCADE;

CREATE UNIQUE INDEX uniq_schedule_date
ON tasks(schedule_id, date)
WHERE schedule_id IS NOT NULL;