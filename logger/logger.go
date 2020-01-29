package logger

import (
	"github.com/greenac/gologger"
)

var goLogger *gologger.GoLogger

func Setup(logPath string) {
	if goLogger == nil {
		goLogger = &gologger.GoLogger{LogLevel: gologger.OutputNormal, LogPath: logPath}
		(*goLogger).Setup()
	} else {
		goLogger.LogPath = logPath
	}
}

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

func Debug(a ...interface{}) {
	setup()
	goLogger.Debug(a)
}
