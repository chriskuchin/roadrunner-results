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
	"github.com/chriskuchin/roadrunner-results/pkg/db"
	"github.com/chriskuchin/roadrunner-results/pkg/google"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Routes(app *firebase.App, db db.DB, assetsFolder string, debug bool) chi.Router {
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("auth client failure")
	}

	r := chi.NewRouter()
	r.Use(middleware.AuthenticationMiddleware(authClient.VerifyIDToken, []string{http.MethodGet, http.MethodOptions}, os.Getenv))

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
				r.Route("/me", func(r chi.Router) {
					r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleGetCurrentUser(db, authClient), "GetCurrentUser"))
				})
				r.Route("/races", func(r chi.Router) {
					r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleRacesList(db), "HandleRacesList"))
					r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleRacesCreate(db), "HandleRacesCreate"))
					r.Route("/import", func(r chi.Router) {
						r.Use(google.HandleOAuth2Creds)
						r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleRaceImportSheet(db), "HandleRaceImportSheet"))
					})
					r.Route("/{raceID}", func(r chi.Router) {
						r.Use(middleware.PathArgCtx("raceID", util.RaceID))
						r.Use(middleware.UserAuthMiddleware(db, []string{http.MethodGet, http.MethodOptions}))

						r.Method(http.MethodDelete, "/", otelhttp.NewHandler(v1.HandleRaceDelete(db), "HandleRaceDelete"))
						r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleRaceGet(db), "HandleRaceGet"))

						r.Route("/stats", func(r chi.Router) {
							r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleRaceStatsGet(db), "HandleRaceStatsGet"))
						})

						r.Route("/volunteers", func(r chi.Router) {
							r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleVolunteersList(db, app), "HandleVolunteersList"))
							r.Method(http.MethodPut, "/", otelhttp.NewHandler(v1.HandleVolunteersCreate(db, app), "HandleVolunteersCreate"))
						})

						r.Route("/participants", func(r chi.Router) {
							r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleParticipantsList(db), "HandleParticipantsList"))
							r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleParticipantsCreate(db), "HandleParticipantsCreate"))
							r.Method(http.MethodGet, "/next_bib", otelhttp.NewHandler(v1.HandleParticipantsNextBibNumber(), "HandleParticipantsNextBibNumber"))
							r.Method(http.MethodGet, "/teams", otelhttp.NewHandler(v1.HandleParticipantsDistinctTeams(), "HandleParticipantsDistinctTeams"))
							r.Method(http.MethodPost, "/csv", otelhttp.NewHandler(v1.HandleParticipantsImportCSV(db), "HandleParticipantsImportCSV"))

							r.Route("/bib/{bibNumber}", func(r chi.Router) {
								r.Use(middleware.PathArgCtx("bibNumber", util.BibNumber))
								r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleParticipantGetByBibNumber(db), "HandleParticipantGetByBibNumber"))
							})
							r.Route("/{participantID}", func(r chi.Router) {
								r.Use(middleware.PathArgCtx("participantID", util.ParticipantID))
								r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleParticipantGet(), "HandleParticipantGet"))
								r.Method(http.MethodPut, "/", otelhttp.NewHandler(v1.HandleParticipantUpdate(db), "HandleParticipantUpdate"))
							})
						})

						r.Route("/events", func(r chi.Router) {
							r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleEventsList(db), "HandleEventsList"))
							r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleEventsCreate(db), "HandleEventsCreate"))

							r.Route("/{eventID}", func(r chi.Router) {
								r.Use(middleware.PathArgCtx("eventID", util.EventID))
								r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleEventGet(db), "HandleEventGet"))
								r.Method(http.MethodDelete, "/", otelhttp.NewHandler(v1.HandleEventDelete(db), "HandleEventDelete"))

								r.Route("/attempts", func(r chi.Router) {
									r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleEventAttemptsCreate(db), "HandleEventAttemptsCreate"))

									r.Route("/{bibNumber}", func(r chi.Router) {
										r.Use(middleware.PathArgCtx("bibNumber", util.BibNumber))
										r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleEventAttemptsList(db), "HandleEventAttemptsList"))
									})
								})
								r.Route("/results", func(r chi.Router) {
									r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleEventResultsGet(db), "HandleEventResultsGet"))
									r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleEventResultsCreate(db, s3, "photo-finish"), "HandleEventResultsCreate"))
									r.Method(http.MethodPut, "/", otelhttp.NewHandler(v1.HandleEventResultsCreate(db, s3, "photo-finish"), "HandleEventResultsCreate"))

									r.Route("/{resultID}", func(r chi.Router) {
										r.Use(middleware.PathArgCtx("resultID", util.ResultID))
										r.Method(http.MethodPatch, "/", otelhttp.NewHandler(v1.HandleEventResultsUpdate(db), "HandleEventResultsUpdate"))
										r.Method(http.MethodDelete, "/", otelhttp.NewHandler(v1.HandleEventResultsDelete(db), "HandleEventResultsDelete"))
									})
								})

								r.Route("/heats", func(r chi.Router) {
									r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleHeatsList(db), "HandleHeatsList"))
									r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleHeatsCreate(db), "HandleHeatsCreate"))

									r.Route("/{timerID}", func(r chi.Router) {
										r.Use(middleware.PathArgCtx("timerID", util.TimerID))
										r.Method(http.MethodPut, "/", otelhttp.NewHandler(v1.HandleHeatUpdate(db), "HandleHeatUpdate"))
										r.Method(http.MethodDelete, "/", otelhttp.NewHandler(v1.HandleHeatDelete(db), "HandleHeatDelete"))
									})
								})

								r.Route("/timers", func(r chi.Router) {
									r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleTimersCreate(db), "HandleTimersCreate"))
									r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleTimersList(db), "HandleTimersList"))

									r.Route("/{timerID}", func(r chi.Router) {
										r.Use(middleware.PathArgCtx("timerID", util.TimerID))
										r.Method(http.MethodPut, "/", otelhttp.NewHandler(v1.HandleTimerStart(db), "HandleTimerStart"))
										r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleTimerGet(db), "HandleTimerGet"))
										r.Method(http.MethodDelete, "/", otelhttp.NewHandler(v1.HandleTimerDelete(db), "HandleTimerDelete"))
									})
								})
							})
						})

						r.Route("/divisions", func(r chi.Router) {
							r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleDivisionsList(db), "HandleDivisionsList"))
							r.Method(http.MethodPost, "/", otelhttp.NewHandler(v1.HandleDivisionsCreate(db), "HandleDivisionsCreate"))
						})
					})
				})

				r.Route("/athletes/results", func(r chi.Router) {
					r.Method(http.MethodGet, "/", otelhttp.NewHandler(v1.HandleAthleteResultsSearch(db), "HandleAthleteResultsSearch"))
				})
			})
		})

		r.Route("/google", func(r chi.Router) {
			r.HandleFunc("/oauth2/callback", v1.HandleGoogleOAuth())
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
