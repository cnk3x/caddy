package errors // import "github.com/pteich/errors"

import (
	"errors"
	"fmt"
)

// New returns a new error from a given string.
func New(msg string) error {
	return errors.New(msg)
}

// Wrap returns an error that wraps an existing error with a new one for a given string.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf returns an error that wraps an existing error with a new one for a given format and arguments.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	args = append(args, err)
	return fmt.Errorf(format+": %w", args...)
}

// Errorf returns an error for a given format and arguments.
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Is checks if a given errors wraps another target.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}
