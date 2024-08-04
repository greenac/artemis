package errs

import (
	"strings"
)

type IGenError interface {
	error
	AddMsg(msg string) IGenError
}

func NewGenError(msg string) IGenError {
	return GenError{Messages: []string{msg}}
}

type GenError struct {
	Messages []string
}

func (e GenError) Error() string {
	bldr := strings.Builder{}
	for i, m := range e.Messages {
		if i == len(e.Messages)-1 {
			bldr.WriteString(m)
		} else {
			bldr.WriteString(m + "->")
		}
	}

	return bldr.String()
}

func (e GenError) AddMsg(msg string) IGenError {
	e.Messages = append(e.Messages, msg)
	return e
}
