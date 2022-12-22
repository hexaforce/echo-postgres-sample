package main

import (
	"echo-postgres-sample/app/api"
	postgres "echo-postgres-sample/app/db"
	"flag"
	"log"

	_ "echo-postgres-sample/app/docs"
	websocket "echo-postgres-sample/app/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	swagger "github.com/swaggo/echo-swagger"
)

var address = flag.String("address", ":1323", "http service address")

func main() {
	flag.Parse()

	hub := websocket.NewHub()
	go hub.Run()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Swagger UI
	e.GET("/swagger/*", swagger.WrapHandler)

	// Websocket lissten
	e.GET("/ws/:userName", func(c echo.Context) error {
		return websocket.ServeWs(hub, c)
	})

	// DB Migration
	pgdb, err := postgres.MigrateDB()
	if err != nil {
		log.Printf("error: %v", err)
		panic("error starting the database")
	}

	// Attach Handler
	api.HandlerMapping(e, pgdb)

	e.Logger.Fatal(e.Start(*address))

}
