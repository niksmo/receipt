package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/niksmo/receipt/pkg/logger"
)

const (
	readHeaderTimeout = 500 * time.Millisecond
	readTimeout       = 2 * time.Second
	idleTimeout       = 1 * time.Second
	handlerTimeout    = 5 * time.Second
	handlerTimeoutMsg = "service unavailable"
)

type wrapper struct {
	Status int
	http.ResponseWriter
}

func (w *wrapper) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type httpServer struct {
	log logger.Logger
	srv *http.Server
	mux *http.ServeMux
}

func New(log logger.Logger, addr string) *httpServer {
	srv := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       2 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
	mux := http.NewServeMux()
	server := &httpServer{log, srv, mux}
	srv.Handler = http.TimeoutHandler(
		server.logResponse(mux), handlerTimeout, handlerTimeoutMsg,
	)
	return server
}

func (s *httpServer) Mux() *http.ServeMux {
	return s.mux
}

func (s *httpServer) Run(ctx context.Context, onFall func(err error)) {
	const op = "httpServer.Run"
	log := s.log.WithOp(op)

	log.Info().Str("addr", s.srv.Addr).Msg("http server is running")
	err := s.srv.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		onFall(fmt.Errorf("%s: %w", op, err))
	}
}

func (s *httpServer) Close() {
	const op = "httpServer.Close"
	log := s.log.WithOp(op)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Info().Msg("start closing server")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to close server gracefully")
		return
	}
	log.Info().Msg("server closed")
}

func (s *httpServer) logResponse(next http.Handler) http.Handler {
	log := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappee := &wrapper{ResponseWriter: w}
		next.ServeHTTP(wrappee, r)
		resDur := time.Since(start)
		resDurStrMicro := strconv.FormatInt(resDur.Microseconds(), 10) + "μs"

		s.log.Info().Str(
			"repliedWithin", resDurStrMicro).Int("status", wrappee.Status).Send()
	}
	return http.HandlerFunc(log)
}
