package golog

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Logger interface {
	Write(msg LogMessage) error
}

type JSONFormatLogger struct {
	Out io.Writer
}

func (j JSONFormatLogger) Write(msg LogMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(j.Out, "%s", b)

	return err
}

type TextFormatLogger struct {
	Stdout, Stderr io.Writer
}

func (t TextFormatLogger) Write(msg LogMessage) error {
	var buf []string

	fixedKey := []string{
		KeyTime,
		KeyLevel,
		KeyMessage,
		KeyError,
	}

	var o io.Writer
	var level Level

	lvl, ok := msg[KeyLevel]
	if !ok {
		return fmt.Errorf("level key not found in LogMessage")
	}

	if err := level.UnmarshalText([]byte(lvl.(string))); err != nil {
		o = t.Stdout
	} else {
		if level == Fatal || level == Error {
			o = t.Stderr
		} else {
			o = t.Stdout
		}
	}

	for _, k := range fixedKey {
		if v, ok := msg[k]; ok {
			switch v.(type) {
			case int:
				buf = append(buf, fmt.Sprintf("%s=%d", k, v))
			default:
				buf = append(buf, fmt.Sprintf("%s=%q", k, v))
			}
			delete(msg, k)
		}
	}

	for k, v := range msg {
		switch v.(type) {
		case int:
			buf = append(buf, fmt.Sprintf("%s=%d", k, v))
		default:
			buf = append(buf, fmt.Sprintf("%s=%q", k, v))
		}
	}

	_, err := fmt.Fprintln(o, strings.Join(buf, " "))

	return err
}
