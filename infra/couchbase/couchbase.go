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

	bucket := cluster.Bucket("dummys")
	bucket.WaitUntilReady(3*time.Second, &gocb.WaitUntilReadyOptions{})

	return &CouchbaseRepository{
		cluster: cluster,
		bucket:  bucket,
		tracer:  tracer,
	}
}

func (r *CouchbaseRepository) GetDummy(ctx context.Context, id string) (*domain.Dummy, error) {
	ctx, span := r.tracer.Wrapped().Start(ctx, "GetDummy")
	defer span.End()

	data, err := r.bucket.DefaultCollection().Get(id, &gocb.GetOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
		ParentSpan: gocbopentelemetry.NewOpenTelemetryRequestSpan(ctx, span),
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("dummy not found")
		}

		zap.L().Error("Failed to get dummy", zap.Error(err))
		return nil, err
	}

	var dummy domain.Dummy
	if err := data.Content(&dummy); err != nil {
		zap.L().Error("Failed to unmarshal dummy", zap.Error(err))
		return nil, err
	}

	return &dummy, nil
}

func (r *CouchbaseRepository) CreateDummy(ctx context.Context, dummy *domain.Dummy) error {
	_, err := r.bucket.DefaultCollection().Insert(dummy.ID, dummy, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})
	if err != nil {
		zap.L().Error("Failed to create dummy", zap.Error(err))
		return err
	}

	return nil
}

func (r *CouchbaseRepository) UpdateDummy(ctx context.Context, dummy *domain.Dummy) error {
	ctx, span := r.tracer.Wrapped().Start(ctx, "UpdateDummy")
	_, err := r.bucket.DefaultCollection().Replace(dummy.ID, dummy, &gocb.ReplaceOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
		ParentSpan: gocbopentelemetry.NewOpenTelemetryRequestSpan(ctx, span),
	})
	if err != nil {
		zap.L().Error("Failed to update dummy", zap.Error(err))
		return err
	}

	return nil
}
