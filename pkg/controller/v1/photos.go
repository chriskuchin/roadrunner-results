package v1

// func (api *Handler) listPhotoFinishFiles(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	output, err := api.s3.ListObjects(ctx, &s3.ListObjectsInput{
// 		Bucket:  aws.String(api.bucket),
// 		Prefix:  aws.String(fmt.Sprintf("/races/%s/events/%s/finish/", util.GetRaceIDFromContext(ctx), util.GetEventIDFromContext(ctx))),
// 		MaxKeys: 500,
// 	})

// 	if err != nil {
// 		w.WriteHeader(http.StatusBadGateway)
// 		return
// 	}

// 	photos := []string{}
// 	for _, v := range output.Contents {
// 		photos = append(photos, url.QueryEscape(*v.Key))
// 		log.Info().Interface("obj", v).Send()
// 	}
// 	render.JSON(w, r, photos)
// }

// func (api *Handler) renderPhotoFinishFile(w http.ResponseWriter, r *http.Request) {
// 	key, _ := url.QueryUnescape(chi.URLParam(r, "photoKey"))
// 	res, _ := api.s3.GetObject(r.Context(), &s3.GetObjectInput{
// 		Bucket: aws.String(api.bucket),
// 		Key:    aws.String(key),
// 	})

// 	img, _ := io.ReadAll(res.Body)
// 	r.Header.Add("Content-Type", "image/png")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(img)
// }
