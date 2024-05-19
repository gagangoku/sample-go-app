package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// Flags
var portFlag = flag.Int("port", 4001, "The port to run server on")
var prettyLogFlag = flag.Bool("prettyLog", true, "Whether to log pretty or json")

type App struct {
}

func (app *App) Init(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msgf("Init called")
}

func (app *App) SetupHttpServer(ctx context.Context) {
	mux := app.Configure(ctx)
	StartServer(ctx, mux, *portFlag)
}

func (app *App) Configure(ctx context.Context) *http.ServeMux {
	logger := zerolog.Ctx(ctx)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/version", helperFn("/version", app.versionHandler))
	mux.HandleFunc("/healthz", helperFn("/healthz", app.healthHandler))
	mux.HandleFunc("/", helperFn("/", app.handleAll))

	logger.Info().Msgf("Server starting")
	return mux
}

func StartServer(baseCtx context.Context, mux *http.ServeMux, port int) {
	logger := zerolog.Ctx(baseCtx)

	ctx, cancelCtx := context.WithCancel(baseCtx)
	serverOne := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, ctxKey, l.Addr().String())
			return ctx
		},
	}

	// Start the server
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Error().Msgf("Error: server one closed")
		} else if err != nil {
			logger.Error().Msgf("Error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	WaitForKillSignal(ctx)
}

func helperFn(handler string, handlerFn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		promHttpLatency.MustCurryWith(prometheus.Labels{"handler": handler}),
		promhttp.InstrumentHandlerCounter(
			promHttpCalls.MustCurryWith(prometheus.Labels{"handler": handler}),
			promhttp.InstrumentHandlerResponseSize(
				promResponseSize,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					newCtx, logger := NewCtxWithUid(r.Context())
					logger.Info().Msgf("[REQ] %s", handler)
					handlerFn(w, r.WithContext(newCtx))
				}),
			),
		),
	)
}

func NewCtxWithUid(ctx context.Context) (context.Context, zerolog.Logger) {
	newUuid := uuid.Must(uuid.NewRandom()).String()
	logger := zerolog.Ctx(ctx).With().Str("u", newUuid).Logger()
	newCtx := logger.WithContext(ctx)
	return newCtx, logger
}

func WaitForKillSignal(ctx context.Context) {
	logger := zerolog.Ctx(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if ctx != nil {
		select {
		case <-ctx.Done():
			logger.Info().Msgf("Context done, exiting")
			return
		case <-c:
			logger.Info().Msgf("Interrupt received, exiting")
			return
		}
	} else {
		for range c {
			logger.Info().Msgf("Interrupt received, exiting")
			return
		}
	}
}

func (app *App) versionHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, VERSION)
}

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func (app *App) handleAll(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello!")
}

type key int

const ctxKey key = iota
