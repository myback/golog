package golog

import (
	"fmt"
	"os"
	"time"
)

const (
	KeyTime    = "time"
	KeyLevel   = "level"
	KeyMessage = "message"
	KeyError   = "error"
)

type Log struct {
	timeFormat string
	level      Level
	logger     Logger
}

type LogMessage map[string]interface{}

func (lm LogMessage) Put(k string, v interface{}) {
	lm[k] = v
}

func (lm LogMessage) PutMessage(v interface{}) {
	lm.Put(KeyMessage, v)
}

func (lm LogMessage) PutError(v interface{}) {
	lm.Put(KeyError, v)
}

func (lm LogMessage) Get(k string, fallback interface{}) interface{} {
	if v, ok := lm[k]; ok {
		return v
	}

	return fallback
}

func New(lvl Level, logger Logger, fmtTime string) *Log {
	if lvl == Custom {
		lvl = Info
	}

	if logger == nil {
		logger = TextFormatLogger{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	}

	if fmtTime == "" {
		fmtTime = "2006-01-02T15:04:05.000Z"
	}

	return &Log{
		timeFormat: fmtTime,
		level:      lvl,
		logger:     logger,
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

func (l *Log) Fatal(msg string) {
	l.Fatalf(LogMessage{KeyMessage: msg})
}

func (l *Log) Fatalf(msg LogMessage) {
	l.print(Fatal, msg)
	os.Exit(127)
}

func (l *Log) Error(msg string) {
	l.Errorf(LogMessage{KeyMessage: msg})
}

func (l *Log) Errorf(msg LogMessage) {
	l.print(Error, msg)
}

func (l *Log) Warning(msg string) {
	l.Warningf(LogMessage{KeyMessage: msg})
}

func (l *Log) Warningf(msg LogMessage) {
	l.print(Warning, msg)
}

func (l *Log) Info(msg string) {
	l.Infof(LogMessage{KeyMessage: msg})
}

func (l *Log) Infof(msg LogMessage) {
	l.print(Info, msg)
}

func (l *Log) Debug(msg string) {
	l.Debugf(LogMessage{KeyMessage: msg})
}

func (l *Log) Debugf(msg LogMessage) {
	l.print(Debug, msg)
}

func (l *Log) Trace(msg string) {
	l.Tracef(LogMessage{KeyMessage: msg})
}

func (l *Log) Tracef(msg LogMessage) {
	l.print(Trace, msg)
}

func (l *Log) print(lvl Level, msg LogMessage) {
	if !l.isLevel(lvl) {
		return
	}

	msg[KeyTime] = time.Now().UTC().Format(l.timeFormat)
	msg[KeyLevel] = lvl.ToString()

	if err := l.logger.Write(msg); err != nil {
		fmt.Fprintf(os.Stderr, "%s ERROR: Print formatter failed: %s", msg["time"], err)
	}
}

func (l Log) NewLogMessage(level string) LogMessage {
	return LogMessage{
		KeyTime:  time.Now().UTC().Format(l.timeFormat),
		KeyLevel: level,
	}
}

func (l *Log) Write(msg LogMessage) {
	lvl, ok := msg.Get(KeyLevel, "").(Level)
	if !ok {
		lvl = Custom
	}

	if l.isLevel(lvl) || lvl == Custom {
		if err := l.logger.Write(msg); err != nil {
			fmt.Fprintf(os.Stderr, "%s ERROR: Print formatter failed: %s", msg["time"], err)
		}
	}
}
