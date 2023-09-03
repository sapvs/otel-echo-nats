package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

func main() {
	cancelTrace, _ := otl.Init(common.FINAL_ROUTE)
	defer cancelTrace(context.Background(), 20*time.Second)

	router := echo.New()
	router.Use(otelecho.Middleware(common.SERVICE_NAME))
	router.GET(common.FINAL_ROUTE, handler)

	log.Fatal(router.Start(":8080"))
}

func handler(c echo.Context) error {
	// get trace from req context
	_, span := otel.Tracer(common.SERVICE_NAME).Start(c.Request().Context(), "final http")
	defer span.End()

	// lets random delay
	sleep := time.Now().UnixMicro() % 500
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	return c.String(http.StatusOK, "great done")

}
