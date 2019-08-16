package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anothrNick/json-tree-service/database"
	"github.com/anothrNick/json-tree-service/websockets"
	"github.com/gin-gonic/gin"
)

// Database is the required interface to the DB layer from the HTTP handlers
type Database interface {
	TranslateError(err error) *database.TranslatedError

	CreateProject(projectName string, data []byte) error
	DeleteProject(projectName string) error

	GetProjectKey(projectName string, keys ...string) ([]byte, error)
	CreateProjectKey(projectName string, data []byte, keys ...string) error
	UpdateProjectKey(projectName string, data []byte, keys ...string) error
	DeleteProjectKey(projectName string, keys ...string) error
}

// Dispatcher provides an interface to dispatch events to clients connect over websockets
type Dispatcher interface {
	Broadcast() chan *websockets.Message
}

// Action is the payload dispatched to any clients connected over websocket
type Action struct {
	Type    string      `json:"type"`
	Project string      `json:"project"`
	Path    string      `json:"path"`
	Data    interface{} `json:"data"`
}

// Handlers contains all handler functions
type Handlers struct {
	db         Database
	dispatcher Dispatcher
}

// NewHandlers creates and returns a new instance of `Handlers` with the datastore
func NewHandlers(datastore Database, dispatcher Dispatcher) *Handlers {
	return &Handlers{
		db:         datastore,
		dispatcher: dispatcher,
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

	// dispatch action
	h.broadcast("POST", keys, project, b)

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

	// dispatch action
	h.broadcast("PUT", keys, project, b)

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

	// dispatch action
	h.broadcast("DELETE", keys, project, nil)

	c.JSON(http.StatusCreated, gin.H{})
}

func (h *Handlers) broadcast(typ, path, project string, b []byte) {
	// dispatch action
	var rawData interface{}
	json.Unmarshal(b, &rawData)
	action := Action{
		Type:    typ,
		Path:    path,
		Project: project,
		Data:    rawData,
	}
	h.dispatcher.Broadcast() <- &websockets.Message{Channel: project, Data: action}
}
