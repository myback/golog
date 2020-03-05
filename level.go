package golog

import (
	"strconv"
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
	levelStr := strings.ToLower(string(lvl))

	if i, err := strconv.Atoi(levelStr); err == nil {
		if Trace < Level(i) {
			*l = Trace
			return nil
		}
	}

	switch levelStr {
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
		*l = Info
	}

	return nil
}
