package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handlers contains all handler functions
type Handlers struct {
	db Database
}

// NewHandlers creates and returns a new instance of `Handlers` with the datastore
func NewHandlers(datastore Database) *Handlers {
	return &Handlers{
		db: datastore,
	}
}

// CreateProject creates a new JSON tree for a project name
func (h *Handlers) CreateProject(c *gin.Context) {
	project := c.Param("project")
	b, _ := c.GetRawData()

	err := h.db.CreateProject(project, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// ReadProjectKey retrieves the data stored at the key path provided by the HTTP path parameters
func (h *Handlers) ReadProjectKey(c *gin.Context) {
	project := c.Param("project")
	keys := c.Param("keys")

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	byt, err := h.db.GetProjectKey(project, strings.Split(keys, "/")...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	var obj interface{}
	json.Unmarshal(byt, &obj)
	c.IndentedJSON(http.StatusOK, obj)
}
