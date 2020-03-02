package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Init struct {
	Level
	Stdout, Stderr io.Writer
	Formatter
}

func (i *Init) New() *Logger {
	stdout := i.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := i.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	return &Logger{
		level:  i.Level,
		stdout: stdout,
		stderr: stderr,
		fmt:    i.Formatter,
	}
}

type Level int

const (
	Fatal Level = iota
	Error
	Warning
	Info
	Debug
)

const Message = "message"

type Logger struct {
	level          Level
	stdout, stderr io.Writer
	fmt            Formatter
}

type LogMessage map[string]interface{}

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
	}

	return ""
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
		return fmt.Errorf("unknown log level %s", lvl)
	}

	return nil
}

func (l *Logger) GetStdout() io.Writer {
	return l.stdout
}

func (l *Logger) GetStderr() io.Writer {
	return l.stderr
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) IsDebug() bool {
	return l.level == Debug
}

func (l *Logger) checkLevel(level Level) bool {
	return l.level >= level
}

func (l *Logger) Errorf(msg LogMessage) {
	l.Log(Error, msg)
}

func (l *Logger) Error(msg string) {
	l.Log(Error, LogMessage{Message: msg})
}

func (l *Logger) Infof(msg LogMessage) {
	l.Log(Info, msg)
}

func (l *Logger) Info(msg string) {
	l.Log(Info, LogMessage{Message: msg})
}

func (l *Logger) Debugf(msg LogMessage) {
	l.Log(Debug, msg)
}

func (l *Logger) Debug(msg string) {
	l.Log(Debug, LogMessage{Message: msg})
}

func (l *Logger) Fatalf(msg LogMessage) {
	l.Log(Fatal, msg)
	os.Exit(127)
}

func (l *Logger) Fatal(msg string) {
	l.Log(Fatal, LogMessage{Message: msg})
	os.Exit(127)
}

func (l *Logger) Log(level Level, msg LogMessage) {
	if l.checkLevel(level) {
		l.CustomLevelLog(level.ToString(), msg)
	}
}

func (l *Logger) CustomLevelLog(level string, msg LogMessage) {
	msg["time"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	msg["level"] = level

	b, err := l.fmt.Format(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: Log formatter failed: %s", msg["time"], err)
		return
	}

	var out io.Writer
	var lvl Level

	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		out = l.stdout
	} else {
		if lvl > Error {
			out = l.stdout
		} else {
			out = l.stderr
		}
	}

	if _, err := fmt.Fprintln(out, b); err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: Write log failed: %s", msg["time"], err)
	}
}
