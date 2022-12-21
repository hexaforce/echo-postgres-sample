package main

import (
	"echo-postgres-sample/app/api"
	"echo-postgres-sample/app/db"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	websocket "echo-postgres-sample/app/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	swagger "github.com/swaggo/echo-swagger"
)

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

	// Websocket
	e.GET("/ws/:userName", func(c echo.Context) error {
		return websocket.ServeWs(hub, c)
	})

	log.Print("server has started")

	//start the db
	pgdb, err := db.StartDB()
	if err != nil {
		log.Printf("error: %v", err)
		panic("error starting the database")
	}

	//get the router of the API by passing the db
	router := api.StartAPI(pgdb)

	//get the port from the environment variable
	port := os.Getenv("PORT")

	//pass the router and start listening with the server
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
		return
	}

}
