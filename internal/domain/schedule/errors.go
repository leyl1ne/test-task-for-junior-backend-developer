package schedule

import "errors"

var (
	ErrInvalidInterval         = errors.New("invalid interval")
	ErrInvalidDayOfMonth       = errors.New("ivalid day of month")
	ErrSpecificDatesIsRequired = errors.New("specific dates is required")
	ErrEvenOddIsRequired       = errors.New("even_odd is required")
)
