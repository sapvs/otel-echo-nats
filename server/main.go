package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sapvs/otelechonats/common"
	"github.com/sapvs/otelechonats/otl"
	"github.com/sapvs/otelechonats/pub"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	server := echo.New()

	cancelTrace, _ := otl.Init("server")
	defer cancelTrace(context.Background(), 20*time.Second)
	// To trace http calls and echo internal traces
	server.Use(otelecho.Middleware(common.SERVICE_NAME))

	server.GET("/hello", getAndNatsPublish)
	server.Logger.Fatal(server.Start(":8080"))
}

func getAndNatsPublish(c echo.Context) error {
	if err := pub.PublishSomeEvent(c.Request().Context()); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	return c.String(http.StatusOK, "Hello\n")
}
