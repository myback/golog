package logging

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	l.write("error", fmt.Sprintf(format, msg...))
}

func (l *Logger) Error(msg interface{}) {
	l.write("error", msgWrap(msg))
}

func (l *Logger) Infof(format string, msg ...interface{}) {
	l.write("info", fmt.Sprintf(format, msg...))
}

func (l *Logger) Info(msg string) {
	l.write("info", msgWrap(msg))
}

func (l *Logger) Debugf(format string, msg ...interface{}) {
	if l.debug {
		l.write("debug", fmt.Sprintf(format, msg...))
	}
}

func (l *Logger) Debug(msg string) {
	if l.debug {
		l.write("debug", msgWrap(msg))
	}
}

func (l *Logger) Fatalf(format string, msg ...interface{}) {
	l.write("fatal", fmt.Sprintf(format, msg...))
	os.Exit(127)
}

func (l *Logger) Fatal(msg string) {
	l.write("fatal", msgWrap(msg))
	os.Exit(127)
}

func (l *Logger) Access(next http.Handler, logHeader bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u, h string
		login, _, ok := r.BasicAuth()
		if ok {
			u = login
		} else {
			u = "-"
		}

		if logHeader {
			var hl []string
			for k, v := range r.Header {
				hl = append(hl, fmt.Sprintf(`http_%s="%s"`, k, v[0]))
			}

			h = " " + strings.Join(hl, " ")
		}

		l.write("access", fmt.Sprintf(`user="%s" address="%s" host="%s" uri="%s" method="%s" proto="%s" useragent="%s" referer="%s"`+h,
			u, r.RemoteAddr, r.Host, r.RequestURI, r.Method, r.Proto, r.UserAgent(), r.Referer()))

		next.ServeHTTP(w, r)
	})
}

func (l *Logger) write(level, msg string) {
	_, err := fmt.Fprintf(l.output, `time="%s" level="%s" %s\n`, time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), level, msg)
	if err != nil {
		fmt.Printf("Log write failed: %s", err)
	}

}

func msgWrap(msg interface{}) string {
	return fmt.Sprintf(`message="%s"`, msg)
}
