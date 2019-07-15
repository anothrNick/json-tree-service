package main

import (
	"log"

	"github.com/anothrNick/json-tree-service/backend/database"
	"github.com/anothrNick/json-tree-service/backend/web"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewPostgres("testuser", "1234", "localhost", "testdb")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	web.SetRoutes(router, db)

	router.Run(":5000")
}
