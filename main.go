package main

import (
	"log"
	"os"

	"github.com/anothrNick/json-tree-service/database"
	"github.com/anothrNick/json-tree-service/web"
	"github.com/anothrNick/json-tree-service/websockets"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
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
	web.SetRoutes(router, httpHandler)
	websockets.SetRoutes(router, wsHandler)

	// run server
	router.Run(":5001")
}
