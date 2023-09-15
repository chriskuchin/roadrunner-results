package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/chriskuchin/roadrunner-results/pkg/services"
	"github.com/chriskuchin/roadrunner-results/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type RecordResultRequestPayload struct {
	Timer       string `json:"timer_id"`
	End         int64  `json:"end_ts"`
	Bib         string `json:"bib_number"`
	ElapsedTime int64  `json:"elapsed_time"`
}

func EventResultsRoutes(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.handleEventResults)
	r.Post("/", handler.recordResult)
	r.Put("/", handler.recordResult)

	r.Get("/photos", handler.listPhotoFinishFiles)
	r.Get("/photos/{photoKey}", handler.renderPhotoFinishFile)

	return r
}

func (api *Handler) listPhotoFinishFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	output, err := api.s3.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket:  aws.String(api.bucket),
		Prefix:  aws.String(fmt.Sprintf("/races/%s/events/%s/finish/", util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))),
		MaxKeys: 500,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	photos := []string{}
	for _, v := range output.Contents {
		photos = append(photos, url.QueryEscape(*v.Key))
		log.Info().Interface("obj", v).Send()
	}
	render.JSON(w, r, photos)
}

func (api *Handler) renderPhotoFinishFile(w http.ResponseWriter, r *http.Request) {
	key, _ := url.QueryUnescape(chi.URLParam(r, "photoKey"))
	res, _ := api.s3.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(api.bucket),
		Key:    aws.String(key),
	})

	img, _ := io.ReadAll(res.Body)
	r.Header.Add("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(img)
}

// gender, team, timer, name
func (api *Handler) handleEventResults(w http.ResponseWriter, r *http.Request) {
	var results []services.ParticipantEventResult
	var filters map[string][]string = map[string][]string{}

	filters["race"] = []string{util.GetRaceIDFromContext(r.Context())}
	filters["event"] = []string{util.GetEventIDFromContext(r.Context())}
	if r.URL.Query().Has("timerId") {
		filters["timer"] = r.URL.Query()["timerId"]
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

	results, err := services.GetEventResults(r.Context(), api.db, filters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, results)
}

func (api *Handler) recordResult(w http.ResponseWriter, r *http.Request) {
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
			_, err = api.s3.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(api.bucket),
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

		ctx := util.SetDB(r.Context(), api.db)
		if payload.End > 0 {
			err = services.RecordTimerResult(ctx, payload.End)
			if err != nil {
				handleBadRequest(err, w, r)
				return
			}
		} else if payload.ElapsedTime > 0 {
			err = services.RecordElapsedTimeResult(ctx, payload.ElapsedTime)
			if err != nil {
				handleBadRequest(err, w, r)
				return
			}
		}

		if payload.Bib != "" {
			timerID := payload.Timer
			if timerID == "" {
				_, timerID, err = services.GetActiveTimerStart(ctx)
				if err != nil {
					handleBadRequest(err, w, r)
					return
				}
			}
			err = services.RecordFinisherResult(ctx, api.db, util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx), timerID, payload.Bib)
			if err != nil {
				handleBadRequest(err, w, r)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
	}

	w.WriteHeader(http.StatusUnsupportedMediaType)
}
