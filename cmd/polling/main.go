package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	pgrepo "kiln-exercice/internal/pg"
	delegationpoll "kiln-exercice/internal/usecase/delegation/poll"
	"kiln-exercice/pkg/pg"
	"kiln-exercice/pkg/tzkt"
)

type Parameters struct {
	DB pg.Parameters

	TzktURL                string    `env:"TZKT_URL" env-default:"https://api.tzkt.io"`
	PollingIntervalSeconds int       `env:"POLLING_INTERVAL_SECONDS" env-default:"10"`
	DefaultPollingFrom     time.Time `env:"DEFAULT_POLLING_FROM" env-layout:"2006-01-02" env-default:"2018-01-01"`
	PollingBatchSize       int       `env:"POLLING_BATCH_SIZE" env-default:"1000"`
}

func main() {
	var (
		ctx    = context.Background()
		params = Parameters{}
	)

	err := env.ReadEnv(&params)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing environment variables")
	}

	db, err := pg.New(params.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}

	tzktSDK, err := tzkt.NewSDK(params.TzktURL)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating tzkt sdk")
	}

	delegationRepo := pgrepo.NewDelegationRepository(db, params.PollingBatchSize)
	pollingRepo := pgrepo.NewPollingRepository(db)

	delegationUseCase := delegationpoll.NewUseCase(
		delegationRepo,
		pollingRepo,
		tzktSDK,
		params.DefaultPollingFrom,
		time.Now,
	)

	ctx, done := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,
	)
	defer done()

	g, ctx := errgroup.WithContext(ctx)

	log.Info().Msg("delegations polling started")

	g.Go(
		func() error {
			ticker := time.NewTicker(time.Duration(params.PollingIntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
					if err = delegationUseCase.PollDelegations(ctx); err != nil {
						return err
					}
				}
			}
		},
	)

	if err = g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("error polling delegations")
	}

	log.Info().Msg("delegations polling stopped")
}
