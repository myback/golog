package golog

import (
	"strings"
)

type Level int

const (
	Custom Level = iota
	Fatal
	Error
	Warning
	Info
	Debug
	Trace
)

func (l Level) ToString() string {
	switch l {
	case Fatal:
		return "fatal"
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	case Debug:
		return "debug"
	case Trace:
		return "trace"
	default:
		return "custom"
	}
}

func (l *Level) UnmarshalText(lvl []byte) error {
	switch strings.ToLower(string(lvl)) {
	case "f", "0", "fatal":
		*l = Fatal
	case "e", "1", "err", "error":
		*l = Error
	case "w", "2", "warn", "warning":
		*l = Warning
	case "i", "3", "info":
		*l = Info
	case "d", "4", "dbg", "debug":
		*l = Debug
	case "t", "5", "trace":
		*l = Trace
	default:
		*l = Custom
	}

	return nil
}
