package task

import "errors"

var (
	ErrNotFound           = errors.New("task not found")
	ErrTaskAlreadyCreated = errors.New("task already created")
)
