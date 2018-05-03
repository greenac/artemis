package logger

import (
	"github.com/greenac/gologger"
)

var goLogger *gologger.GoLogger

func setup() {
	if goLogger == nil {
		goLogger = &gologger.GoLogger{LogLevel: gologger.OutputNormal, LogPath: ""}
		(*goLogger).Setup()
	}
}

func Log(a ...interface{}) {
	setup()
	goLogger.Log(a)
}

func Error(a ...interface{}) {
	setup()
	goLogger.Error(a)
}

func Warn(a ...interface{}) {
	setup()
	goLogger.Warn(a)
}
