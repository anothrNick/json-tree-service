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
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// ReadProject retrieves the project root JSON tree
func (h *Handlers) ReadProject(c *gin.Context) {
	project := c.Param("project")

	byt, err := h.db.GetProjectKey(project)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	var obj interface{}
	json.Unmarshal(byt, &obj)
	c.IndentedJSON(http.StatusOK, obj)
}

// DeleteProject deletes the entire project
func (h *Handlers) DeleteProject(c *gin.Context) {
	project := c.Param("project")

	err := h.db.DeleteProject(project)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ReadProjectKey retrieves the data stored at the key path provided by the HTTP path parameters
func (h *Handlers) ReadProjectKey(c *gin.Context) {
	project := c.Param("project")
	keys := c.Param("keys")

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	byt, err := h.db.GetProjectKey(project, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	var obj interface{}
	json.Unmarshal(byt, &obj)
	c.IndentedJSON(http.StatusOK, obj)
}

// CreateProjectKey updates a project key at the key path. An error is returned if the key already exists.
func (h *Handlers) CreateProjectKey(c *gin.Context) {
	project := c.Param("project")
	keys := c.Param("keys")

	b, _ := c.GetRawData()

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.CreateProjectKey(project, b, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateProjectKey updates a project key at the key path. The key is created if it does not already exist.
func (h *Handlers) UpdateProjectKey(c *gin.Context) {
	project := c.Param("project")
	keys := c.Param("keys")

	b, _ := c.GetRawData()

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.UpdateProjectKey(project, b, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteProjectKey deletes a project key at the key path
func (h *Handlers) DeleteProjectKey(c *gin.Context) {
	project := c.Param("project")
	keys := c.Param("keys")
	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.DeleteProjectKey(project, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
