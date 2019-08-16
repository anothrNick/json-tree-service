package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// OpenCORSMiddleware controls the cross origin policies.
func OpenCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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
	api := engine.Group("/api")

	api.Use(OpenCORSMiddleware())

	api.GET("/:project", h.ReadProject)      // returns entire root tree
	api.POST("/:project", h.CreateProject)   // create a new tree at `project`
	api.DELETE("/:project", h.DeleteProject) // project tree must be empty to delete

	api.GET("/:project/*keys", h.ReadProjectKey)
	api.POST("/:project/*keys", h.CreateProjectKey)
	api.PUT("/:project/*keys", h.UpdateProjectKey)
	api.DELETE("/:project/*keys", h.DeleteProjectKey)

	return nil
}
