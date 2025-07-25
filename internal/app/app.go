package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/niksmo/receipt/internal/httpapi"
	"github.com/niksmo/receipt/internal/service"
	"github.com/niksmo/receipt/pkg/logger"
)

type Opts struct {
	Addr     string
	Login    string
	Password string
	SMTPHost string
	SMTPPort string
}

type App struct {
	log  logger.Logger
	s    *http.Server
	opts Opts
}

func New(log logger.Logger, opts Opts) *App {
	server := &http.Server{
		Addr: opts.Addr,
	}

	app := &App{
		log:  log,
		s:    server,
		opts: opts,
	}

	app.initMux()

	return app
}

func (app *App) Run(done context.CancelFunc) {
	const op = "App.Run"
	log := app.log.WithOp(op)
	log.Info().Msg("start listen and serve")
	err := app.s.ListenAndServe()

	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("unexpected server shutdown")
		}
	}
	log.Info().Msg("stop listen and server")
	done()
}

func (app *App) Close() {
	const op = "App.Close"
	log := app.log.WithOp(op)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("start closing server")
	if err := app.s.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn().Msg("close server deadline exceeded")
			return
		}
		log.Error().Err(err).Msg("unexpected error while closing")
		return
	}
	log.Info().Msg("server closed successfully")
}

func (app *App) initMux() {
	mux := http.NewServeMux()
	sender := service.NewSMTPTransmitter(
		app.opts.Login,
		app.opts.Password,
		app.opts.SMTPHost,
		app.opts.SMTPPort,
	)
	emailNotifier := service.NewEmailNotifier(sender)
	httpapi.RegisterMailReceiptHandler(app.log, mux, emailNotifier)
	app.s.Handler = mux
}
