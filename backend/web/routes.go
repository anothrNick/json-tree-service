package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// SetRoutes sets all of the appropriate routes to handlers for the application
func SetRoutes(engine *gin.Engine, datastore Database) error {

	handlers := NewHandlers(datastore)

	engine.POST("/:project", handlers.CreateProject) // create a new tree at `project`
	engine.GET("/:project", notImplemented)          // returns entire root tree
	engine.DELETE("/:project", notImplemented)       // project tree must be empty to delete

	engine.GET("/:project/*keys", handlers.ReadProjectKey)
	engine.POST("/:project/*keys", notImplemented)
	engine.PUT("/:project/*keys", notImplemented)
	engine.DELETE("/:project/*keys", notImplemented)

	return nil
}
