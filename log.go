package golog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	KeyTime           = "time"
	KeyLevel          = "level"
	KeyMessage        = "message"
	KeyError          = "error"
	defaultTimeFormat = "2006-01-02T15:04:05.000Z"
)

var logger Logger

type Log struct {
	timeFormat string
	level      Level
	//logger     Logger
}

type LogMessage map[string]interface{}

func (lm *LogMessage) Put(k string, v interface{}) *LogMessage {
	(*lm)[k] = v

	return lm
}

func (lm *LogMessage) PutMessage(v interface{}) *LogMessage {
	lm.Put(KeyMessage, v)

	return lm
}

func (lm *LogMessage) PutError(v interface{}) *LogMessage {
	lm.Put(KeyError, v)

	return lm
}

func (lm *LogMessage) AddTrace() *LogMessage {
	pc := make([]uintptr, 5)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var frame runtime.Frame
	for {
		frame, _ = frames.Next()
		if strings.Contains(frame.Function, "golog") {
			if frame.Line == 0 {
				lm.Put("package", "failed")
				break
			}

			continue
		}

		lm.Put("package", frame.Function)
		lm.Put("file", fmt.Sprintf("%s:%d", frame.File, frame.Line))
		break
	}

	return lm
}

func (lm *LogMessage) Get(k string, fallback interface{}) interface{} {
	if v, ok := (*lm)[k]; ok {
		return v
	}

	return fallback
}

func New(lvl Level, l Logger, fmtTime string) *Log {
	if lvl == Custom {
		lvl = Info
	}

	if l == nil {
		logger = TextFormatLogger{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	} else {
		logger = l
	}

	if fmtTime == "" {
		fmtTime = defaultTimeFormat
	}

	return &Log{
		timeFormat: fmtTime,
		level:      lvl,
		//logger:     logger,
	}
}

func (l *Log) GetLevel() Level {
	return l.level
}

func (l *Log) IsDebug() bool {
	return l.level == Debug
}

func (l *Log) isLevel(level Level) bool {
	return l.level >= level
}

func (l *Log) Fatal(msg interface{}) {
	l.Fatalf(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Fatalf(msg LogMessage) {
	l.print(Fatal, msg)
	os.Exit(127)
}

func (l *Log) Error(msg interface{}) {
	l.Errorf(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Errorf(msg LogMessage) {
	l.print(Error, msg)
}

func (l *Log) Warning(msg interface{}) {
	l.Warningf(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Warningf(msg LogMessage) {
	l.print(Warning, msg)
}

func (l *Log) Info(msg interface{}) {
	l.Infof(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Infof(msg LogMessage) {
	l.print(Info, msg)
}

func (l *Log) Debug(msg interface{}) {
	l.Debugf(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Debugf(msg LogMessage) {
	l.print(Debug, msg)
}

func (l *Log) Trace(msg interface{}) {
	l.Tracef(LogMessage{KeyMessage: fmt.Sprintf("%s", msg)})
}

func (l *Log) Tracef(msg LogMessage) {
	msg.AddTrace()
	l.print(Trace, msg)
}

func (l *Log) print(lvl Level, msg LogMessage) {
	if !l.isLevel(lvl) {
		return
	}

	msg[KeyTime] = time.Now().UTC().Format(l.timeFormat)
	msg[KeyLevel] = lvl.ToString()

	if err := logger.Write(msg); err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: Print formatter failed: %s", msg["time"], err)
	}
}

func (l *Log) NewLogMessage(level string) LogMessage {
	return LogMessage{
		KeyTime:  time.Now().UTC().Format(l.timeFormat),
		KeyLevel: level,
	}
}

func SampleLogging(level string) LogMessage {
	return LogMessage{
		KeyTime:  time.Now().UTC().Format(defaultTimeFormat),
		KeyLevel: level,
	}
}

func (lm *LogMessage) Write() {
	l := logger
	if l == nil {
		l = TextFormatLogger{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}

	if err := l.Write(*lm); err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: Print formatter failed: %s", (*lm)[KeyTime], err)
	}
}

func (l *Log) Write(msg LogMessage) {
	if _, ok := msg[KeyTime]; !ok {
		msg[KeyTime] = time.Now().UTC().Format(l.timeFormat)
	}

	lvl := msg.Get(KeyLevel, "custom").(Level)

	if l.isLevel(lvl) || lvl == Custom {
		if err := logger.Write(msg); err != nil {
			fmt.Fprintf(os.Stderr, "%s ERROR: Print formatter failed: %s", msg[KeyTime], err)
		}
	}
}
