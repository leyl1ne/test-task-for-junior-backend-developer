package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
	"example.com/taskservice/internal/transport/http/util"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	Schedule *scheduleDTO `json:"schedule,omitempty"`
	Date     *string      `json:"date,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type scheduleDTO struct {
	Type string `json:"type"`

	Interval      *int     `json:"interval,omitempty"`
	DaysOfMonth   []int    `json:"days_of_month,omitempty"`
	SpecificDates []string `json:"specific_dates,omitempty"`
	EvenOdd       *string  `json:"even_odd,omitempty"`

	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
}

func mapScheduleDTO(dto *scheduleDTO) (*taskusecase.ScheduleInput, error) {
	start, err := util.ParseDate(dto.StartDate)
	if err != nil {
		return nil, err
	}

	var end *time.Time
	if dto.EndDate != nil {
		d, err := util.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, err
		}
		end = &d
	}

	var specific []time.Time
	for _, s := range dto.SpecificDates {
		d, err := util.ParseDate(s)
		if err != nil {
			return nil, err
		}
		specific = append(specific, d)
	}

	return &taskusecase.ScheduleInput{
		Type:          dto.Type,
		Interval:      dto.Interval,
		DaysOfMonth:   dto.DaysOfMonth,
		SpecificDates: specific,
		EvenOdd:       dto.EvenOdd,
		StartDate:     start,
		EndDate:       end,
	}, nil
}
