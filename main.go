package main

import (
	"fmt"
	"net/http"
	"os"

	"database/sql"

	"github.com/chriskuchin/roadrunner-results/pkg/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	port   string
	debug  bool
	dbPath string
)

func main() {
	app := &cli.App{
		Name: "results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Usage:       "Port for the server to listen on",
				Aliases:     []string{"p"},
				EnvVars:     []string{"PORT"},
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
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "launches the application web server",
				Action: func(c *cli.Context) error {
					zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
					// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
					zerolog.SetGlobalLevel(zerolog.InfoLevel)
					if debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}

					sql.Open("sqlite3", dbPath)
					r := chi.NewRouter()

					// A good base middleware stack
					r.Use(middleware.RequestID)
					r.Use(middleware.RealIP)
					r.Use(middleware.Logger)
					r.Use(middleware.Recoverer)

					r.Route("/", func(r chi.Router) {
						r.Mount("/healthcheck", controller.HealthcheckResources{}.Routes())
						r.Mount("/api", controller.APIsResource{}.Routes())

						r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
							http.FileServer(http.Dir("./public/dist")).ServeHTTP(w, r)
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
