package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	scheduledomain "example.com/taskservice/internal/domain/schedule"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	taskRepo     TaskRepository
	scheduleRepo ScheduleRepository
	now          func() time.Time
}

func NewService(taskRepo TaskRepository, scheduleRepo ScheduleRepository) *Service {
	return &Service{
		taskRepo:     taskRepo,
		scheduleRepo: scheduleRepo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()

	if input.Schedule != nil {
		schedule := mapToDomainSchedule(input.Schedule, normalized.Title, now)

		if err := schedule.Validate(); err != nil {
			return nil, err
		}

		schedule, err := s.scheduleRepo.Create(ctx, schedule)
		if err != nil {
			return nil, err
		}

		today := scheduledomain.NormalizeDate(now)

		if schedule.ShouldRunToday(today) {
			task := buildTaskFromSchedule(schedule, today, now)
			return s.taskRepo.Create(ctx, task)
		}

		return nil, nil
	}

	if input.Date == nil {
		return nil, fmt.Errorf("%w: date required", ErrInvalidInput)
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Date:        input.Date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.taskRepo.Create(ctx, model)
}

func buildTaskFromSchedule(s *scheduledomain.Schedule, date time.Time, now time.Time) *taskdomain.Task {
	return &taskdomain.Task{
		Title:      s.Title,
		Status:     taskdomain.StatusNew,
		ScheduleID: &s.ID,
		Date:       &date,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.taskRepo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}

	updated, err := s.taskRepo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.taskRepo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.taskRepo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}
