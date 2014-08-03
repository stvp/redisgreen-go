package redisgreen

import (
	"strings"
)

// Error is a convenience struct for returning RedisGreen JSON API errors as Go
// errors.
type Error struct {
	Errors []string `json:"errors"`
}

func (e *Error) Error() string {
	return strings.Join(e.Errors, " ")
}
