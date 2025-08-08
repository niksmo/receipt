package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/niksmo/receipt/pkg/logger"
)

func AcceptJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.Header().Set("Accept", "application/json")
			errStr := "invalid media type"
			http.Error(w, errStr, http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogResposeStatus(l logger.Logger, next http.Handler) http.Handler {
	return httpLog{l, next}
}

type httpLog struct {
	log  logger.Logger
	next http.Handler
}

type logWriteHeaderWrapper struct {
	statusCode int
	http.ResponseWriter
}

func (whw *logWriteHeaderWrapper) WriteHeader(statusCode int) {
	whw.statusCode = statusCode
	whw.ResponseWriter.WriteHeader(statusCode)
}

func (l httpLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wrapper := &logWriteHeaderWrapper{ResponseWriter: w}
	l.next.ServeHTTP(wrapper, r)
	procDur := time.Since(start)
	resDurStrMicro := strconv.FormatInt(procDur.Microseconds(), 10) + "μs"
	statusCode := wrapper.statusCode

	l.log.Info().Str(
		"procDur", resDurStrMicro).Int("status", statusCode).Send()
}
