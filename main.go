package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://testuser:1234@localhost/testdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/:project", func(c *gin.Context) {
		project := c.Param("project")
		b, _ := c.GetRawData()

		_, err = db.Exec("INSERT INTO trees (project, data) VALUES ($1, $2)", project, b)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, gin.H{})
	})

	// router.GET("/:project", func(c *gin.Context) {
	// 	//project := c.Param("project")
	// })

	router.GET("/:project/*keys", func(c *gin.Context) {
		project := c.Param("project")
		keys := c.Param("keys")
		byt := []byte{}

		keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
		keys = strings.Join(strings.Split(keys, "/"), ",")
		err = db.QueryRow(fmt.Sprintf("SELECT data#>'{%s}' as data FROM trees WHERE project=$1 ORDER BY id DESC LIMIT 1", keys), project).Scan(&byt)
		if err != nil {
			log.Fatal(err)
		}

		var obj interface{}
		json.Unmarshal(byt, &obj)
		c.IndentedJSON(http.StatusOK, obj)
	})

	router.Run(":5000")
}
