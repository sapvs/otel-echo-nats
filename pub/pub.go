package pub

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(common.NATS_URL)
}

func PublishSomeEvent(ctx context.Context) error {
	ctx, span := otel.GetTracerProvider().
		Tracer(common.SERVICE_NAME).
		Start(ctx, "nats_publish", trace.WithSpanKind(trace.SpanKindProducer))

	defer span.End()

	msg := nats.NewMsg(common.NATS_SUBJECT)
	log.Printf("msg.Header before: %v\n", msg.Header)
	otl.InjectOtelToNATSHeader(ctx, &msg.Header)
	log.Printf("msg.Header After: %v\n", msg.Header)

	span.SetAttributes(attribute.String("header", fmt.Sprintf("%v", msg.Header)))

	if err := nc.PublishMsg(msg); err != nil {
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "error in publishing")
		return err
	}
	return nil
}
