package task

import (
	"context"
	"time"

	scheduledomain "example.com/taskservice/internal/domain/schedule"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type TaskRepository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type ScheduleRepository interface {
	Create(ctx context.Context, s *scheduledomain.Schedule) (*scheduledomain.Schedule, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type CreateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	Date        *time.Time
	Schedule    *ScheduleInput
}

type UpdateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
}

type ScheduleInput struct {
	Type string

	Interval      *int
	DaysOfMonth   []int
	SpecificDates []time.Time
	EvenOdd       *string

	StartDate time.Time
	EndDate   *time.Time
}

func mapToDomainSchedule(input *ScheduleInput, title string, now time.Time) *scheduledomain.Schedule {
	return &scheduledomain.Schedule{
		Title:         title,
		Type:          scheduledomain.Type(input.Type),
		Interval:      input.Interval,
		DaysOfMonth:   input.DaysOfMonth,
		SpecificDates: input.SpecificDates,
		EvenOdd:       input.EvenOdd,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
