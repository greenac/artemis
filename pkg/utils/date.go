package utils

import (
	"time"
)

const DateFormat = "2006-01-02T15:04:05.999Z"

func CurrentIsoTime() string {
	return time.Now().UTC().Format(DateFormat)
}
