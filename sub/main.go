package main

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"go.opentelemetry.io/otel"
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
	defer sub.Unsubscribe()

	for {
		if msg, err := sub.NextMsg(10 * time.Second); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Message received %v :  %s", msg.Header, msg.Data)
			otlCtx := otl.ExtractOTelContextFromNATSHeader(context.Background(), &msg.Header)
			_, span := otel.Tracer(common.SERVICE_NAME).Start(otlCtx, "subject receive", trace.WithSpanKind(trace.SpanKindConsumer))

			span.End()
		}

	}
}
