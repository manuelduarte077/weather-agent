package api

import (
	"net/http"
	"strings"
)

// ValidateCity validates and sanitizes the city parameter.
// Returns the sanitized city and an error if validation fails.
func ValidateCity(city string, maxLength int) (string, error) {
	if city == "" {
		return "", &ValidationError{
			Message:    "city parameter cannot be empty",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Sanitize city input
	city = strings.TrimSpace(city)
	if city == "" {
		return "", &ValidationError{
			Message:    "Invalid city parameter",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Limit city length to prevent abuse
	if maxLength > 0 && len(city) > maxLength {
		return "", &ValidationError{
			Message:    "City parameter too long",
			StatusCode: http.StatusBadRequest,
		}
	}

	return city, nil
}

// ValidationError represents a validation error with HTTP status code.
type ValidationError struct {
	Message    string
	StatusCode int
}

func (e *ValidationError) Error() string {
	return e.Message
}

