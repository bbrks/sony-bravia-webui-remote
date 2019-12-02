package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/zchee/color"
)

type LogFn func(level LogLevel, ctx context.Context, format string, v ...interface{})

type Logger interface {
	Log(level LogLevel, ctx context.Context, format string, v ...interface{})
}

type LogLevel uint8

const (
	LevelAll   = LogLevel(iota)
	LevelError = LogLevel(iota)
	LevelWarn
	LevelInfo
	LevelDebug
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "[DBG]"
	case LevelInfo:
		return "[INF]"
	case LevelWarn:
		return "[WRN]"
	case LevelError:
		return "[ERR]"
	case LevelAll:
		return "[ALL]"
	}
	return ""
}

var _ Logger = &defaultLogger{}

type defaultLogger struct {
	w          io.Writer
	timeFormat string
	color      bool
	quiet      bool
}

func newDefaultLogger(quiet bool) *defaultLogger {
	shouldUseColor := os.Getenv("NO_COLOR") == "" && !color.NoColor

	return &defaultLogger{
		w:          os.Stderr,
		timeFormat: "2006-01-02T15:04:05.000Z07:00", // RFC3339 Milli
		color:      shouldUseColor,
		quiet:      quiet,
	}
}

func (l *defaultLogger) Log(logLevel LogLevel, ctx context.Context, format string, v ...interface{}) {
	if logLevel >= LevelWarn && l.quiet {
		return
	}
	_, _ = fmt.Fprintf(l.w, l.formatLog(ctx, logLevel, format), v...)
}

// formatLog adds the common prefixes to log format strings
func (l *defaultLogger) formatLog(ctx context.Context, level LogLevel, format string) string {
	timestampPrefix := time.Now().Format(l.timeFormat) + " "

	levelPrefix := level.String() + " "

	contextPrefix := ""
	if ctx != nil {
		if reqID, ok := ctx.Value(reqIDKey{}).(string); ok {
			contextPrefix = "req:" + reqID + " "
		}
	}

	colorFn := noColor
	if l.color {
		colorFn = level.colorFn()
	}

	return timestampPrefix + colorFn(levelPrefix+contextPrefix+format) + "\n"
}

// colorFn returns a Sprintf function that colors a string for the given log level.
func (l LogLevel) colorFn() (colorFunc func(string, ...interface{}) string) {
	switch l {
	case LevelDebug:
		return color.HiBlackString
	case LevelWarn:
		return color.YellowString
	case LevelError:
		return color.RedString
	}
	// no color
	return noColor
}

func noColor(s string, v ...interface{}) string {
	return s
}
