package v1

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	s3     *s3.Client
	db     *sqlx.DB
	app    *firebase.App
	debug  bool
	bucket string
}

func Routes(db *sqlx.DB, app *firebase.App, debug bool) chi.Router {
	r := chi.NewRouter()

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

	handler := &Handler{
		db:     db,
		debug:  debug,
		app:    app,
		s3:     s3.NewFromConfig(cfg),
		bucket: "photo-finish",
	}

	r.Mount("/races", RacesRoutes(handler))
	r.Mount("/google", GoogleRoutes(handler))
	return r
}
