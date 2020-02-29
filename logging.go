package logging

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	debug  bool
	output io.Writer
}

func New(out io.Writer, debug bool) *Logger {
	return &Logger{
		debug:  debug,
		output: out,
	}
}

func (l *Logger) Errorf(format string, msg ...interface{}) {
	_print("error", _msgWrap(format, msg...))
}

func (l *Logger) Infof(format string, msg ...interface{}) {
	_print("info", _msgWrap(format, msg...))
}

func (l *Logger) Debugf(format string, msg ...interface{}) {
	if l.debug {
		_print("debug", _msgWrap(format, msg...))
	}
}

func (l *Logger) Fatalf(format string, msg ...interface{}) {
	_print("fatal", _msgWrap(format, msg...))
	os.Exit(127)
}

func _print(level, msg string) {
	fmt.Printf(`time="%s" level="%s" %s`, time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), level, msg)
	fmt.Println()
}

func _msgWrap(format string, msg ...interface{}) string {
	if len(msg) == 0 {
		return fmt.Sprintf(`message="%s"`, format)
	}

	return fmt.Sprintf(format, msg...)
}
