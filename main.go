package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chriskuchin/roadrunner-results/pkg/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
						zerolog.SetGlobalLevel(zerolog.InfoLevel)
						log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
					}

					log.Debug().Str("db", dbPath).Send()
					db, err := sqlx.Open("sqlite3", dbPath)
					if err != nil {
						log.Fatal().Err(err).Send()
					}

					r := chi.NewRouter()
					// A good base middleware stack
					r.Use(middleware.RequestID)
					r.Use(middleware.RealIP)
					r.Use(middleware.Logger)
					r.Use(middleware.Recoverer)

					r.Route("/", func(r chi.Router) {
						r.Mount("/healthcheck", controller.HealthcheckResources{}.Routes())
						r.Mount("/api", controller.Routes(db, debug))

						r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
							http.FileServer(http.Dir(frontendFolder)).ServeHTTP(w, r)
						})
					})

					log.Debug().Msgf("Launching server listening on: %s", port)
					http.ListenAndServe(fmt.Sprintf(":%s", port), r)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
