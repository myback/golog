package logging

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Formatter interface {
	Format(msg LogMessage) (string, error)
}

type JSONFormatter LogMessage
type TextFormatter LogMessage

func (j JSONFormatter) Format(msg LogMessage) (string, error) {
	b, err := json.Marshal(msg)
	return string(b), err
}

func (t TextFormatter) Format(msg LogMessage) (string, error) {
	var buf []string

	fixedKey := []string{
		"time",
		"level",
		"message",
	}

	for _, k := range fixedKey {
		if v, ok := msg[k]; ok {
			buf = append(buf, fmt.Sprintf("%s=\"%s\"", k, v))
			delete(msg, k)
		}
	}

	for k, v := range msg {
		buf = append(buf, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	return strings.Join(buf, " "), nil
}
