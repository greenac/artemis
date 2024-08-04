package gologger

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"reflect"
	"time"
)

const ShowDebugKey = "GO_LOGGER_SHOW_DEBUG"

type outputType int

const (
	OutputError   outputType = 1000
	OutputDebug   outputType = 1001
	OutputNormal  outputType = 1002
	OutputWarning outputType = 1003
)

func innerElement(a []interface{}) []interface{} {
	aa := a[0]
	v := reflect.ValueOf(aa)
	if v.Kind() != reflect.Slice {
		return a
	}

	return innerElement(aa.([]interface{}))
}

type GoLogger struct {
	LogLevel   outputType
	LogPath    string
	timeFormat string
	isSetup    bool
	showDebug  bool
}

// TODO: Add support for different log levels
func (l *GoLogger) Setup() {
	if !l.isSetup {
		if l.timeFormat == "" {
			l.timeFormat = time.UnixDate
		}

		l.showDebug = os.Getenv(ShowDebugKey) == "true"
		l.isSetup = true
	}
}

func (l *GoLogger) coloredOutput(ot outputType, a ...interface{}) {
	var c *color.Color
	switch ot {
	case OutputDebug:
		c = color.New(color.FgHiMagenta).Add(color.Bold)
	case OutputError:
		c = color.New(color.FgRed).Add(color.Bold)
	case OutputWarning:
		c = color.New(color.FgYellow).Add(color.Bold)
	default:
		c = color.New(color.FgCyan)
	}

	c.Println(a...)
}

func (l *GoLogger) writeToFile(message string) {
	l.Setup()
	if l.LogPath == "" {
		return
	}

	go func(msg string) {
		f, err := os.OpenFile(l.LogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			l.coloredOutput(OutputError, "Error: could not write to octopus log file:", l.LogPath)
			return
		}
		defer f.Close()

		msg += "\n"
		_, err = f.WriteString(msg)
		if err != nil {
			l.coloredOutput(OutputError, "Error: failed to write message to log file:", l.LogPath)
		}
	}(message)
}

func (l *GoLogger) log(ot outputType, a ...interface{}) {
	l.Setup()
	aa := innerElement(a)
	args := fmt.Sprint(aa)
	var pre string
	switch ot {
	case OutputError:
		pre = "ERROR: "
	case OutputWarning:
		pre = "WARNING: "
	case OutputNormal:
		pre = "LOG: "
	case OutputDebug:
		pre = "DEBUG: "
	default:
		fmt.Println("Error: output type:", ot, "is unknown")
		pre = ""
	}

	msg := time.Now().Format(time.UnixDate) + " " + pre + args[1: len(args) - 1]
	l.coloredOutput(ot, msg)
	l.writeToFile(msg)
}

func (l *GoLogger) Log(a ...interface{}) {
	l.log(OutputNormal, a)
}

func (l *GoLogger) Error(a ...interface{}) {
	l.log(OutputError, a)
}

func (l *GoLogger) Debug(a ...interface{}) {
	if l.showDebug {
		l.log(OutputDebug, a)
	}
}

func (l *GoLogger) Warn(a ...interface{}) {
	l.log(OutputWarning, a)
}
