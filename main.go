package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/chriskuchin/roadrunner-results/pkg/controller"
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	port           string
	debug          bool
	dbPath         string
	frontendFolder string
)

func main() {
	app := &cli.App{
		Name: "results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Usage:       "Port for the server to listen on",
				Aliases:     []string{"p"},
				EnvVars:     []string{"PORT", "NOMAD_HOST_PORT_run"},
				Destination: &port,
				Value:       "3030",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable Debug Logs",
				Aliases:     []string{"verbose"},
				EnvVars:     []string{""},
				Destination: &debug,
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "db-path",
				Usage:       "Specifies the path of the database",
				Aliases:     []string{"db"},
				EnvVars:     []string{"DB_PATH"},
				Destination: &dbPath,
				Value:       "./results.db",
			},
			&cli.StringFlag{
				Name:        "frontend-folder",
				Usage:       "The frontend folder for the web assets",
				Aliases:     []string{"dist"},
				EnvVars:     []string{"FRONTEND_FOLDER"},
				Destination: &frontendFolder,
				Value:       "./dist",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "launches the application web server",
				Action: func(c *cli.Context) error {
					zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
					zerolog.SetGlobalLevel(zerolog.InfoLevel)
					if debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
						log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
					}

					db := db.NewSQLite(dbPath)

					app, err := firebase.NewApp(context.Background(), nil)
					if err != nil {
						log.Fatal().Msgf("error initializing app: %v\n", err)
					}

					r := chi.NewRouter()
					// A good base middleware stack
					r.Use(middleware.RequestID)
					r.Use(middleware.RealIP)
					r.Use(middleware.Logger)
					r.Use(middleware.Recoverer)

					r.Mount("/", controller.Routes(app, db, frontendFolder, debug))

					bsp := honeycomb.NewBaggageSpanProcessor()

					otelshutdown, err := otelconfig.ConfigureOpenTelemetry(
						otelconfig.WithSpanProcessor(bsp),
					)
					if err != nil {
						log.Fatal().Err(err).Send()
					}

					defer otelshutdown()

					log.Debug().Msgf("Launching server listening on: %s", port)
					err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
					if err != nil {
						log.Error().Err(err).Send()
					}

					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
