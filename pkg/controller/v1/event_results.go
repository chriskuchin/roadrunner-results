package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"

	apiutil "github.com/chriskuchin/roadrunner-results/pkg/controller/api-util"
	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func HandleEventResultsUpdate(db *sqlx.DB) http.HandlerFunc {
	type PatchResultPayload struct {
		Time int    `json:"result"`
		Bib  string `json:"bib_number"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		var payload PatchResultPayload
		err = json.Unmarshal(body, &payload)
		if err != nil {
			apiutil.HandleBadRequest(err, w, r)
			return
		}

		if payload.Bib != "" {
			log.Info().Msg("Bib Update")
		}

		if payload.Time != 0 {
			log.Info().Msg("Time Update")
		}
	}
}

func HandleEventResultsGet(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var results []services.ParticipantEventResult
		var filters map[string][]string = map[string][]string{}

		filters["race"] = []string{util.GetRaceIDFromContext(r.Context())}
		filters["event"] = []string{util.GetEventIDFromContext(r.Context())}
		if r.URL.Query().Has("timerId") {
			if r.URL.Query()["timerId"][0] == "latest" {
				_, timerId, err := services.GetActiveTimerStart(r.Context(), db)
				if err != nil {
					apiutil.HandleBadRequest(err, w, r)
					return
				}

				filters["timer"] = []string{timerId}
			} else {
				filters["timer"] = r.URL.Query()["timerId"]
			}
		}

		if r.URL.Query().Has("gender") {
			filters["gender"] = r.URL.Query()["gender"]
		}

		if r.URL.Query().Has("name") {
			filters["name"] = r.URL.Query()["name"]
		}

		if r.URL.Query().Has("team") {
			filters["team"] = r.URL.Query()["team"]
		}

		if r.URL.Query().Has("year") {
			filters["year"] = r.URL.Query()["year"]
		}

		results, err := services.GetEventResults(r.Context(), db, filters)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, results)
	}
}

func HandleEventResultsCreate(db *sqlx.DB, s3Client *s3.Client, bucket string) http.HandlerFunc {
	type RecordResultRequestPayload struct {
		Timer       string `json:"timer_id"`
		End         int64  `json:"end_ts"`
		Bib         string `json:"bib_number"`
		ElapsedTime int64  `json:"elapsed_time"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "multipart/form-data") {
			reader, err := r.MultipartReader()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				b, err := io.ReadAll(part)
				if err != nil {
					log.Error().Err(err).Send()
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var filename = part.FileName()
				_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(fmt.Sprintf("/races/%s/events/%s/finish/%s", util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), filename)),
					Body:   bytes.NewReader(b),
				})

				if err != nil {
					log.Error().Err(err).Send()
					w.WriteHeader(http.StatusBadGateway)
					return
				}
			}

			w.WriteHeader(http.StatusAccepted)
			return
		} else {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Err(err).Send()
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var payload RecordResultRequestPayload
			err = json.Unmarshal(body, &payload)
			if err != nil {
				log.Error().Err(err).Send()
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ctx := util.SetDB(r.Context(), db)
			if payload.End > 0 {
				err = services.RecordTimerResult(ctx, payload.End)
				if err != nil {
					apiutil.HandleBadRequest(err, w, r)
					return
				}
			} else if payload.ElapsedTime > 0 {
				err = services.RecordElapsedTimeResult(ctx, payload.ElapsedTime)
				if err != nil {
					apiutil.HandleBadRequest(err, w, r)
					return
				}
			}

			if payload.Bib != "" {
				timerID := payload.Timer
				if timerID == "" {
					_, timerID, err = services.GetActiveTimerStart(ctx, db)
					if err != nil {
						apiutil.HandleBadRequest(err, w, r)
						return
					}
				}
				err = services.RecordFinisherResult(ctx, db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, payload.Bib)
				if err != nil {
					apiutil.HandleBadRequest(err, w, r)
					return
				}
			}

			w.WriteHeader(http.StatusAccepted)
		}

		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}
