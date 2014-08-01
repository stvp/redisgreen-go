package redisgreen

import (
	"strings"
)

type Error struct {
	Errors []string `json:"errors"`
}

func (e *Error) Error() string {
	return strings.Join(e.Errors, " ")
}
