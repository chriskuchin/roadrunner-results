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

func (api *Handler) handleEventResults(w http.ResponseWriter, r *http.Request) {
	var results []services.ParticipantEventResult
	var err error

	timerID := r.URL.Query().Get("timerId")

	if timerID != "" {
		results, err = services.GetEventHeatResults(r.Context(), api.db, timerID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		results, err = services.GetEventResults(r.Context(), api.db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
			services.RecordTimerResult(ctx, payload.End)
		} else if payload.ElapsedTime > 0 {
			services.RecordElapsedTimeResult(ctx, payload.ElapsedTime)
		}

		if payload.Bib != "" {
			services.RecordFinisherResult(ctx, api.db, payload.Bib)
		}

		w.WriteHeader(http.StatusAccepted)
	}

	w.WriteHeader(http.StatusUnsupportedMediaType)
}
