package main

import (
	"log"
	"os"

	"github.com/anothrNick/json-tree-service/database"
	"github.com/anothrNick/json-tree-service/web"
	"github.com/anothrNick/json-tree-service/websockets"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	// Making change to test workflows

	// connect to db
	db, err := database.NewPostgres(
		os.Getenv("JSONDB_USER"),
		os.Getenv("JSONDB_PW"),
		os.Getenv("JSONDB_HOST"),
		os.Getenv("JSONDB_DB"),
	)

	if err != nil {
		log.Fatal(err)
	}

	// create a new message dispatcher for websocket connections
	dispatcher := websockets.NewDispatcher()
	go dispatcher.Run()

	// create HTTP handlers
	httpHandler := web.NewHandlers(db, dispatcher)

	// create websocket handlers
	wsHandler := websockets.NewHandlers(dispatcher)

	// create router, set routes
	router := gin.Default()

	// serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))

	// serve HTTP routes
	web.SetRoutes(router, httpHandler)

	// serve Websocket routes
	websockets.SetRoutes(router, wsHandler)

	// run server
	router.Run(":5001")
}
