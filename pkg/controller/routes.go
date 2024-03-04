package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chriskuchin/roadrunner-results/pkg/controller/middleware"
	v1 "github.com/chriskuchin/roadrunner-results/pkg/controller/v1"
	"github.com/chriskuchin/roadrunner-results/pkg/google"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func Routes(app *firebase.App, db *sqlx.DB, assetsFolder string, debug bool) chi.Router {
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("auth client failure")
	}

	r := chi.NewRouter()
	r.Use(middleware.AuthenticationMiddleware(authClient))

	if debug {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Api-Token"},
			MaxAge:         300,
		}))
	}

	var accountId = os.Getenv("R2_ACCOUNT_ID")
	var accessKeyId = os.Getenv("R2_ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("R2_SECRET_ACCESS_KEY")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	s3 := s3.NewFromConfig(cfg)

	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/healthcheck", HandleHealthcheckGet())
			r.Route("/v1", func(r chi.Router) {
				r.Route("/races", func(r chi.Router) {
					r.Get("/", v1.HandleRacesList(db))
					r.Post("/", v1.HandleRacesCreate(db))
					r.Route("/import", func(r chi.Router) {
						r.Use(google.HandleOAuth2Creds)
						r.Get("/", v1.HandleRaceImportSheet(db))
					})
					r.Route("/{raceID}", func(r chi.Router) {
						r.Use(middleware.RaceCtx)
						r.Use(middleware.UserAuthMiddleware(db))

						r.Delete("/", v1.HandleRaceDelete(db))
						r.Get("/", v1.HandleRaceGet(db))
						r.Route("/volunteers", func(r chi.Router) {
							r.Put("/", v1.HandleVolunteersCreate(db, app))
							r.Get("/", v1.HandleVolunteersList(db, app))
						})
						r.Route("/participants", func(r chi.Router) {
							r.Get("/", v1.HandleParticipantsList(db))
							r.Post("/", v1.HandleParticipantsCreate(db))
							r.Get("/next_bib", v1.HandleParticipantsNextBibNumber())
							r.Get("/teams", v1.HandleParticipantsDistinctTeams())
							r.Post("/csv", v1.HandleParticipantsImportCSV(db))
							r.Route("/bib/{bib_number}", func(r chi.Router) {
								r.Use(middleware.BibNumberCtx)
								r.Get("/", v1.HandleParticipantGetByBibNumber(db))
							})
							r.Route("/{participantID}", func(r chi.Router) {
								r.Use(middleware.ParticipantCtx)
								r.Get("/", v1.HandleParticipantGet())
								r.Put("/", v1.HandleParticipantUpdate(db))
							})
						})
						r.Route("/events", func(r chi.Router) {
							r.Get("/", v1.HandleEventsList(db))
							r.Post("/", v1.HandleEventsCreate(db))

							r.Route("/{eventID}", func(r chi.Router) {
								r.Use(middleware.EventCtx)
								r.Get("/", v1.HandleEventGet(db))
								r.Delete("/", v1.HandleEventDelete(db))
								r.Route("/attempts", func(r chi.Router) {
									r.Post("/", v1.HandleEventAttemptsCreate(db))
									r.Route("/{bib_number}", func(r chi.Router) {
										r.Use(middleware.BibNumberCtx)
										r.Get("/", v1.HandleEventAttemptsList(db))
									})
								})
								r.Route("/results", func(r chi.Router) {
									r.Get("/", v1.HandleEventResultsGet(db))
									r.Post("/", v1.HandleEventResultsCreate(db, s3, "photo-finish"))
									r.Put("/", v1.HandleEventResultsCreate(db, s3, "photo-finish"))
									r.Route("/{resultID}", func(r chi.Router) {
										r.Use(middleware.ResultCtx)
										r.Patch("/", v1.HandleEventResultsUpdate(db))
									})
								})
								r.Route("/timers", func(r chi.Router) {
									r.Post("/", v1.HandleTimersCreate(db))
									r.Get("/", v1.HandleTimersList(db))

									r.Route("/{timerID}", func(r chi.Router) {
										r.Use(middleware.TimerCtx)
										r.Put("/", v1.HandleTimerStart(db))
										r.Get("/", v1.HandleTimerGet(db))
										r.Delete("/", v1.HandleTimerDelete(db))
									})
								})
							})
						})
						r.Route("/divisions", func(r chi.Router) {
							r.Get("/", v1.HandleDivisionsList(db))
							r.Post("/", v1.HandleDivisionsCreate(db))
						})
					})
				})
			})
		})

		r.Route("/google", func(r chi.Router) {
			r.HandleFunc("/oauth2/callback", v1.HandleGoogleOAuth())
		})
		r.Route("/athletes/results", func(r chi.Router) {
			r.Get("/search", v1.HandleAthleteResultsSearch(db))
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			_, err := os.Stat(filepath.Join(assetsFolder, r.URL.Path))
			if os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(assetsFolder, "index.html"))
			} else {
				http.FileServer(http.Dir(assetsFolder)).ServeHTTP(w, r)
			}
		})
	})

	return r
}
