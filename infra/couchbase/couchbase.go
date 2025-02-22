package couchbase

import (
	"context"
	"errors"
	"time"

	gocbopentelemetry "github.com/couchbase/gocb-opentelemetry"
	"github.com/couchbase/gocb/v2"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"go-template/domain"
)

type CouchbaseRepository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
	tp      *sdktrace.TracerProvider
	tracer  *gocbopentelemetry.OpenTelemetryRequestTracer
}

func NewCouchbaseRepository(tp *sdktrace.TracerProvider) *CouchbaseRepository {
	tracer := gocbopentelemetry.NewOpenTelemetryRequestTracer(tp)
	cluster, err := gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout: 3 * time.Second,
			KVTimeout:      3 * time.Second,
			QueryTimeout:   3 * time.Second,
		},
		Authenticator: gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "123456789",
		},
		Transcoder: gocb.NewJSONTranscoder(),
		Tracer:     tracer,
	})
	if err != nil {
		zap.L().Fatal("Failed to connect to couchbase", zap.Error(err))
	}

	bucket := cluster.Bucket("products")
	bucket.WaitUntilReady(3*time.Second, &gocb.WaitUntilReadyOptions{})

	return &CouchbaseRepository{
		cluster: cluster,
		bucket:  bucket,
		tracer:  tracer,
	}
}

func (r *CouchbaseRepository) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	ctx, span := r.tracer.Wrapped().Start(ctx, "GetProduct")
	defer span.End()

	data, err := r.bucket.DefaultCollection().Get(id, &gocb.GetOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
		ParentSpan: gocbopentelemetry.NewOpenTelemetryRequestSpan(ctx, span),
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("product not found")
		}

		zap.L().Error("Failed to get product", zap.Error(err))
		return nil, err
	}

	var product domain.Product
	if err := data.Content(&product); err != nil {
		zap.L().Error("Failed to unmarshal product", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *CouchbaseRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.bucket.DefaultCollection().Insert(product.ID, product, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})
	if err != nil {
		zap.L().Error("Failed to create product", zap.Error(err))
		return err
	}

	return nil
}

func (r *CouchbaseRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := r.tracer.Wrapped().Start(ctx, "UpdateProduct")
	_, err := r.bucket.DefaultCollection().Replace(product.ID, product, &gocb.ReplaceOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
		ParentSpan: gocbopentelemetry.NewOpenTelemetryRequestSpan(ctx, span),
	})
	if err != nil {
		zap.L().Error("Failed to update product", zap.Error(err))
		return err
	}

	return nil
}
