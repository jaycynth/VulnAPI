package helper

import (
	"errors"
	"fmt"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	ErrInvalidData      = errors.New("invalid data provided")
	ErrUserNotFound     = errors.New("not found")
	ErrMissingData      = errors.New("missing fields")
	ErrDuplicateData    = errors.New("the record already exists")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

func ErrorFormatter(err error, values ...interface{}) error {
	switch {
	case errors.Is(err, ErrInvalidData):
		return fmt.Errorf("error: %w", fmt.Errorf("invalid data provided"))
	case errors.Is(err, ErrUserNotFound):
		return fmt.Errorf("error: %w", fmt.Errorf("user not found"))
	case errors.Is(err, ErrMissingData):
		return fmt.Errorf("error: %w", fmt.Errorf("missing data: %v", values))
	case errors.Is(err, ErrDuplicateData):
		return fmt.Errorf("error: %w", fmt.Errorf("duplicate data entry"))
	case errors.Is(err, ErrMethodNotAllowed):
		return fmt.Errorf("error: %w", fmt.Errorf("method not allowed"))
	default:
		return fmt.Errorf("error: %w", err)
	}
}
