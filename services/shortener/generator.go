package shortener

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const defaultMaxRetries = 10

// ShortURLGenerator generates unique short codes for shortened URLs.
type ShortURLGenerator struct {
	length     int
	maxRetries int
}

// NewShortURLGenerator creates a new ShortURLGenerator.
// length must be > 0 and not exceed the length of a UUID string (~36 chars).
func NewShortURLGenerator(length int) *ShortURLGenerator {
	return &ShortURLGenerator{
		length:     length,
		maxRetries: defaultMaxRetries,
	}
}

// GenerateUniqueShortCode generates a short code and ensures it doesn't already exist.
func (g *ShortURLGenerator) GenerateUniqueShortCode(checkExists func(string) (bool, error)) (string, error) {
	if g.length <= 0 || g.length > 36 {
		return "", &GeneratorError{
			UserMessage: "Invalid short URL length",
			TechMessage: fmt.Sprintf("length must be between 1 and 36, got %d", g.length),
			Err:         errors.New("invalid length"),
		}
	}

	for attempt := 0; attempt < g.maxRetries; attempt++ {
		code := uuid.New().String()[:g.length]

		exists, err := checkExists(code)
		if err != nil {
			return "", &GeneratorError{
				UserMessage: "Internal server error",
				TechMessage: "failed to check if short code exists",
				Err:         err,
			}
		}
		if !exists {
			return code, nil
		}
	}

	return "", &GeneratorError{
		UserMessage: "Could not generate a unique short code",
		TechMessage: fmt.Sprintf("exceeded %d attempts", g.maxRetries),
		Err:         errors.New("code generation exhausted retries"),
	}
}
