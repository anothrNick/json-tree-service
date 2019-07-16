package main

import (
	"log"
	"os"

	"github.com/anothrNick/json-tree-service/database"
	"github.com/anothrNick/json-tree-service/web"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// connect to db
	db, err := database.NewPostgres(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PW"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_DB"),
	)

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
