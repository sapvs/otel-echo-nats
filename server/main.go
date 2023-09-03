package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	nats "github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(common.NATS_URL)
}
func main() {
	server := echo.New()

	cancelTrace, _ := otl.Init("server")
	defer cancelTrace(context.Background(), 20*time.Second)

	server.Use(otelecho.Middleware(common.SERVICE_NAME))

	server.GET("/hello", getAndNatsPublish)

	server.Logger.Fatal(server.Start(":8080"))

}

func getAndNatsPublish(c echo.Context) error {

	if err := publishEvent(c.Request().Context()); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusOK, "Hello\n")
}

func publishEvent(ctx context.Context) error {
	ctx, span := otel.Tracer(common.SERVICE_NAME).Start(ctx, "pub1", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	msg := nats.NewMsg(common.NATS_SUBJECT)
	fmt.Printf("msg.Header before: %v\n", msg.Header)
	otl.InjectOtelToNATSHeader(ctx, &msg.Header)
	fmt.Printf("msg.Header After: %v\n", msg.Header)

	return nc.PublishMsg(msg)
}
