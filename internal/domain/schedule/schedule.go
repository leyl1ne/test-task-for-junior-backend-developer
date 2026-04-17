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

func (s *Schedule) ShouldRunToday(today time.Time) bool {
	// до старта — не запускаем
	if today.Before(NormalizeDate(s.StartDate)) {
		return false
	}

	// после конца — не запускаем
	if s.EndDate != nil && today.After(NormalizeDate(*s.EndDate)) {
		return false
	}

	switch s.Type {

	case Daily:
		if s.Interval == nil {
			return false
		}

		days := int(today.Sub(NormalizeDate(s.StartDate)).Hours() / 24)
		return days%*s.Interval == 0

	case Monthly:
		for _, d := range s.DaysOfMonth {
			if d == today.Day() {
				return true
			}
		}
		return false

	case Specific:
		for _, d := range s.SpecificDates {
			if NormalizeDate(d).Equal(today) {
				return true
			}
		}
		return false

	case EvenOdd:
		if s.EvenOdd == nil {
			return false
		}

		if *s.EvenOdd == "even" {
			return today.Day()%2 == 0
		}
		return today.Day()%2 != 0
	}

	return false
}

func NormalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
