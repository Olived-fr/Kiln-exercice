package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	httphandler "kiln-exercice/internal/handler/delegation"
	pgrepo "kiln-exercice/internal/pg"
	"kiln-exercice/internal/usecase/delegation/list"
	"kiln-exercice/pkg/pg"
)

var (
	serverPort               = flag.String("port", "8080", "server port")
	gracefulShutdownTimeout  = flag.Duration("graceful-shutdown-timeout", 15*time.Second, "graceful shutdown timeout")
	gracelessShutdownTimeout = flag.Duration("graceless-shutdown-timeout", 15*time.Second, "graceless shutdown timeout")
	readTimeout              = flag.Duration("read-timeout", 15*time.Second, "read timeout")
)

type Parameters struct {
	DB pg.Parameters
}

func main() {
	var (
		ctx    = context.Background()
		params = Parameters{}
	)

	flag.Parse()

	err := env.ReadEnv(&params)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing environment variables")
	}

	db, err := pg.New(params.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}
	defer db.Close()

	delegationRepo := pgrepo.NewDelegationRepository(db, 0) // no insert in the api

	// For the sake of simplicity, we define the routes here.
	r := http.NewServeMux()
	r.Handle("GET /xtz/delegations", httphandler.NewDelegationHandler(list.NewUseCase(delegationRepo)))

	if err = listenAndServe(ctx, r); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}

	log.Info().Msg("Server stopped")
}

func listenAndServe(ctx context.Context, mux *http.ServeMux) error {
	g, ctx := errgroup.WithContext(ctx)

	ctx, cancel := signal.NotifyContext(
		ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,
	) // Listen for interrupt signals
	defer cancel()

	reqCtx, reqCancel := context.WithCancel(context.Background())

	server := http.Server{
		Handler:     mux,
		Addr:        ":" + *serverPort,
		ReadTimeout: *readTimeout,
		BaseContext: func(net.Listener) context.Context {
			return reqCtx // Sets the parent context for each incoming request allowing to cancel all of them.
		},
	}

	g.Go(
		func() error {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		},
	)

	log.Info().Msgf("Server started on port %s", *serverPort)

	<-ctx.Done()

	// The server shutdown process is divided into two steps:
	// 1. Graceful shutdown: Process the remaining requests until gracefulShutdownTimeout.
	// 2. Graceless shutdown: Cancel the remaining requests and shutdown the server after gracelessShutdownTimeout.

	timer := time.AfterFunc(*gracefulShutdownTimeout, reqCancel)
	defer timer.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(), *gracefulShutdownTimeout+*gracelessShutdownTimeout,
	)
	defer shutdownCancel()

	g.Go(
		func() error {
			return server.Shutdown(shutdownCtx)
		},
	)

	return errors.Join(g.Wait())
}
