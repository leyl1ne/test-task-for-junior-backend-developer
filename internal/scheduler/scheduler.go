package scheduler

import (
	"context"
	"errors"
	"log"
	"time"

	scheduledomain "example.com/taskservice/internal/domain/schedule"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type ScheduleRepository interface {
	ListActive(ctx context.Context) ([]scheduledomain.Schedule, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
}

type Scheduler struct {
	scheduleRepo ScheduleRepository
	taskRepo     TaskRepository

	interval time.Duration
	now      func() time.Time
}

func New(scheduleRepo ScheduleRepository, taskRepo TaskRepository, interval time.Duration) *Scheduler {
	return &Scheduler{
		scheduleRepo: scheduleRepo,
		taskRepo:     taskRepo,
		interval:     interval,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	now := s.now()
	today := scheduledomain.NormalizeDate(now)

	schedules, err := s.scheduleRepo.ListActive(ctx)
	if err != nil {
		log.Println("scheduler: failed to load schedules:", err)
		return
	}

	for _, sch := range schedules {

		if !sch.ShouldRunToday(today) {
			continue
		}

		task := buildTaskFromSchedule(&sch, today, now)

		_, err := s.taskRepo.Create(ctx, task)
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}

			log.Println("scheduler: failed to create task:", err)
		}
	}
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

func isUniqueViolation(err error) bool {
	return err != nil && errors.Is(err, taskdomain.ErrTaskAlreadyCreated)
}
