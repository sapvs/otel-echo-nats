package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	nats "github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(common.NATS_URL)
}
func main() {
	server := echo.New()

	server.GET("/hello", getAndNatsPublish)

	server.Logger.Fatal(server.Start(":8080"))

}

func getAndNatsPublish(c echo.Context) error {
	if err := nc.Publish(common.NATS_SUBJECT, []byte(time.Now().String())); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Logger().SetLevel(log.INFO)
	c.Logger().Info("publushed")
	return c.String(http.StatusOK, "Hello\n")
}
