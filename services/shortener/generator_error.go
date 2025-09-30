package shortener

import "fmt"

// GeneratorError provides structured error information.
type GeneratorError struct {
	UserMessage string
	TechMessage string
	Err         error
}

func (e *GeneratorError) Error() string {
	// return technical details for logging
	return fmt.Sprintf("%s: %v", e.TechMessage, e.Err)
}

func (e *GeneratorError) Unwrap() error {
	return e.Err
}
