package schedule

import "time"

type Type string

const (
	Daily    Type = "daily"
	Monthly  Type = "monthly"
	Specific Type = "specific"
	EvenOdd  Type = "even_odd"
)

type Schedule struct {
	ID            int64
	Title         string
	Type          Type
	Interval      *int
	DaysOfMonth   []int
	SpecificDates []time.Time
	EvenOdd       *string
	StartDate     time.Time
	EndDate       *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s *Schedule) Validate() error {
	switch s.Type {
	case Daily:
		if s.Interval == nil || *s.Interval <= 0 {
			return ErrInvalidInterval
		}
	case Monthly:
		if !(1 <= len(s.DaysOfMonth) && len(s.DaysOfMonth) <= 30) {
			return ErrInvalidDayOfMonth
		}
	case Specific:
		if len(s.SpecificDates) == 0 {
			return ErrSpecificDatesIsRequired
		}
	case EvenOdd:
		if s.EvenOdd == nil {
			return ErrEvenOddIsRequired
		}
	}

	return nil
}
