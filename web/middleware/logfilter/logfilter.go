// overriding parts of chi's middleware.Logger because we want to use zerolog
package logfilter

import (
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func Setup() {
	middleware.DefaultLogger = middleware.RequestLogger(&zerologLogFormatter{})
}

// --- implement middleware.LogFormatter

type zerologLogFormatter struct {
}

func (l *zerologLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &zerologLogEntry{
		zerologLogFormatter: l,
		request:             r,
	}
	entry.requestId = middleware.GetReqID(r.Context())
	entry.method = r.Method
	entry.path = r.URL.Path
	entry.ip = r.RemoteAddr
	entry.userAgent = r.UserAgent()

	return entry
}

// --- implement middleware.LogEntry

type zerologLogEntry struct {
	*zerologLogFormatter
	request  *http.Request
	requestId string
	method string
	path string
	ip string
	userAgent string
}

func (l *zerologLogEntry) Write(status, bytes int, elapsed time.Duration) {
	msg := "Request"

	var e *zerolog.Event
	switch {
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		e = log.Warn()
	case status >= http.StatusInternalServerError:
		e = log.Error()
	default:
		e = log.Info()
	}

	e.Int("status", status).
		Str("method", l.method).
		Str("path", l.path).
		Str("ip", l.ip).
		Dur("latency", elapsed).
		Str("user-agent", l.userAgent).
		Str("request-id", l.requestId).
		Msg(msg)
}

func (l *zerologLogEntry) Panic(v interface{}, stack []byte) {
	panicEntry := l.NewLogEntry(l.request).(*zerologLogEntry)

	msg := "Request Panic"

	e := log.Panic()

	e.Str("method", panicEntry.method).
		Str("path", panicEntry.path).
		Str("ip", panicEntry.ip).
		Str("user-agent", panicEntry.userAgent).
		Str("request-id", panicEntry.requestId).
		Msg(msg)
}
