package otl

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sapvs/otelechonats/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	prop = otel.GetTextMapPropagator()
)

func InjectOtelToNATSHeader(ctx context.Context, headers *nats.Header) {
	prop.Inject(ctx, propagation.HeaderCarrier(*headers))
	log.Printf(" SAPAN , %v\n", headers)
}

func ExtractOTelContextFromNATSHeader(ctx context.Context, headers *nats.Header) context.Context {
	log.Printf(" SAPAN , %v\n", headers)
	return prop.Extract(ctx, propagation.HeaderCarrier(*headers))
}

func newJaegerProvider(id string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost("localhost"),
			jaeger.WithAgentPort("6831"),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "jaeger.New")
	}

	return tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(common.SERVICE_NAME),
			semconv.DeploymentEnvironmentKey.String("dev"),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.ServiceInstanceIDKey.String(id),
		)),
	), nil
}
func Init(id string) (func(ctx context.Context, shutdownTimeout time.Duration), error) {
	tp, err := newJaegerProvider(id)
	if err != nil {
		return nil, errors.Wrap(err, "newJaegerProvider")
	}

	otel.SetTracerProvider(tp)

	tc := propagation.TraceContext{}
	// Register the TraceContext propagator globally.
	otel.SetTextMapPropagator(tc)

	return func(ctx context.Context, shutdownTimeout time.Duration) {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal("trace provider shutdown", err)
		}
	}, nil
}
