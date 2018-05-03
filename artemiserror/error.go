package artemiserror

import (
	"github.com/greenac/artemis/logger"
)

type ExceptionType string

const (
	ArgsNotInitialized ExceptionType = "ArgsNotInitialized"
)

type ArtemisError struct {
	Message string
}

func (ani *ArtemisError) Error() string {
	return ani.Message
}

func GetArtemisError(et ExceptionType, message *string) *ArtemisError {
	m := ""
	if message == nil {
		switch et {
		case ArgsNotInitialized:
			m = string(ArgsNotInitialized)
		default:
			logger.Warn("No Error of type:", et)
		}
	} else {
		m = *message
	}

	return &ArtemisError{Message: m}
}
