package postgres

import (
	"context"

	scheduledomain "example.com/taskservice/internal/domain/schedule"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Create(ctx context.Context, s *scheduledomain.Schedule) (*scheduledomain.Schedule, error) {
	const query = `
	INSERT INTO schedules (
		title, type, interval, days_of_month, specific_dates,
		even_odd, start_date, end_date, created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		s.Title,
		s.Type,
		s.Interval,
		s.DaysOfMonth,
		s.SpecificDates,
		s.EvenOdd,
		s.StartDate,
		s.EndDate,
		s.CreatedAt,
		s.UpdatedAt,
	).Scan(&s.ID)

	return s, err
}

func (r *ScheduleRepository) ListActive(ctx context.Context) ([]scheduledomain.Schedule, error) {
	const query = `
	SELECT id, title, type, interval, days_of_month,
	       specific_dates, even_odd, start_date, end_date
	FROM schedules
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]scheduledomain.Schedule, 0)
	for rows.Next() {
		schedule, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, *schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

type scheduleScanner interface {
	Scan(dest ...any) error
}

func scanSchedule(scanner scheduleScanner) (*scheduledomain.Schedule, error) {
	var (
		schedule     scheduledomain.Schedule
		typeSchedule string
	)

	if err := scanner.Scan(
		&schedule.ID,
		&schedule.Title,
		&typeSchedule,
		&schedule.Interval,
		&schedule.DaysOfMonth,
		&schedule.SpecificDates,
		&schedule.EvenOdd,
		&schedule.StartDate,
		&schedule.EndDate,
	); err != nil {
		return nil, err
	}

	schedule.Type = scheduledomain.Type(typeSchedule)

	return &schedule, nil
}
