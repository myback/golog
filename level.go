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
)

func (l Level) ToString() string {
	switch l {
	case Fatal:
		return "Fatal"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	default:
		return "Custom"
	}
}

func (l *Level) UnmarshalText(lvl []byte) error {
	switch strings.ToLower(string(lvl)) {
	case "f", "fatal":
		*l = Fatal
	case "e", "err", "error":
		*l = Error
	case "w", "warn", "warning":
		*l = Warning
	case "i", "info":
		*l = Info
	case "d", "dbg", "debug":
		*l = Debug
	default:
		*l = Custom
	}

	return nil
}
