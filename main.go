package main

import (
	"log"

	"github.com/anothrNick/json-tree-service/database"
	"github.com/anothrNick/json-tree-service/web"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// connect to db
	db, err := database.NewPostgres("testuser", "1234", "localhost", "testdb")
	if err != nil {
		log.Fatal(err)
	}

	// create handlers
	handler := web.NewHandlers(db)

	// create router, set routes
	router := gin.Default()
	web.SetRoutes(router, handler)

	// run server
	router.Run(":5000")
}
