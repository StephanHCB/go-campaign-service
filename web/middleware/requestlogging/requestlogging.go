// overriding parts of chi's middleware.Logger because we want to use go-autumn-logging
package requestlogging

import (
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	auloggingapi "github.com/StephanHCB/go-autumn-logging/api"
	"github.com/go-chi/chi/middleware"
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

func (l *zerologLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	msg := "Request"

	ctxLogger := aulogging.Logger.Ctx(l.request.Context())
	var e auloggingapi.LeveledLoggingImplementation
	switch {
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		e = ctxLogger.Warn()
	case status >= http.StatusInternalServerError:
		e = ctxLogger.Error()
	default:
		e = ctxLogger.Info()
	}

	e.With("status", fmt.Sprintf("%d", status)).
		With("method", l.method).
		With("path", l.path).
		With("ip", l.ip).
		With("latency", fmt.Sprintf( "%d", elapsed.Milliseconds())).
		With("user-agent", l.userAgent).
		With("request-id", l.requestId).
		Print(msg)
}

func (l *zerologLogEntry) Panic(v interface{}, stack []byte) {
	panicEntry := l.NewLogEntry(l.request).(*zerologLogEntry)

	msg := "Request Panic"

	e := aulogging.Logger.NoCtx().Panic()

	e.With("method", panicEntry.method).
		With("path", panicEntry.path).
		With("ip", panicEntry.ip).
		With("user-agent", panicEntry.userAgent).
		With("request-id", panicEntry.requestId).
		Print(msg)
}
