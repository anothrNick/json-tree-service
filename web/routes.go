package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// Handler is an interface to the HTTP handler functions.
type Handler interface {
	CreateProject(c *gin.Context)
	ReadProject(c *gin.Context)
	DeleteProject(c *gin.Context)
	ReadProjectKey(c *gin.Context)
	CreateProjectKey(c *gin.Context)
	UpdateProjectKey(c *gin.Context)
	DeleteProjectKey(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for the application
func SetRoutes(engine *gin.Engine, h Handler) error {

	//handlers := NewHandlers(datastore)

	engine.POST("/:project", h.CreateProject)  // create a new tree at `project`
	engine.GET("/:project", h.ReadProject)     // returns entire root tree
	engine.DELETE("/:project", notImplemented) // project tree must be empty to delete

	engine.GET("/:project/*keys", h.ReadProjectKey)
	engine.POST("/:project/*keys", notImplemented)
	engine.PUT("/:project/*keys", notImplemented)
	engine.DELETE("/:project/*keys", notImplemented)

	return nil
}