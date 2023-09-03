package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(common.NATS_URL)
}

func main() {
	cancelTrace, _ := otl.Init("subscriber")
	defer cancelTrace(context.Background(), 20*time.Second)

	sub, er := nc.SubscribeSync(common.NATS_SUBJECT)
	if er != nil {
		log.Fatal(er)
	}
	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("could not unsubscribe due to %s", err.Error())
		}
	}()

	for {
		if msg, err := sub.NextMsg(100 * time.Second); err != nil {
			log.Fatal(err)
		} else {
			otlCtx := otl.ExtractOTelContextFromNATSHeader(context.Background(), &msg.Header)
			_, span := otel.Tracer(common.SERVICE_NAME).Start(otlCtx, "subject receive", trace.WithSpanKind(trace.SpanKindConsumer))

			// call another endpoint
			callHttpFinal(otlCtx)
			span.End()
		}

	}
}

func callHttpFinal(ctx context.Context) {
	callCtx, span := otel.Tracer(common.SERVICE_NAME).Start(ctx, "final http get")
	defer span.End()

	request, err := http.NewRequestWithContext(callCtx, http.MethodGet, "http://last:8080/"+common.FINAL_ROUTE, nil)
	if err != nil {
		log.Println("failed to request final", err.Error())
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println("req timeout")
		span.SetStatus(codes.Error, "failed to request final")
		span.RecordError(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
	}
}

var client = http.Client{
	Timeout: 200 * time.Millisecond,
}
