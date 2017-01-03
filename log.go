package imco

import (
	"log"
	"time"
)

var logger *log.Logger

var trackingStart bool

var trackingT0 time.Time

// Debug show debug message
func Debug(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(format+"\n", args...)
}

// DebugTime show elapsed time with Debug()
func DebugTime(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	if trackingStart {
		trackingStart = false
		args = append(args, time.Since(trackingT0))
		Debug(format+": elapsed %s", args...)
	} else {
		trackingStart = true
		trackingT0 = time.Now()
		Debug(format, args...)
	}
}
