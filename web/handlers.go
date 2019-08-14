package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/anothrNick/json-tree-service/websockets"
	"github.com/gin-gonic/gin"
)

type Action struct {
	Type string      `json:"type"`
	Path string      `json:"path"`
	Data interface{} `json:"data"`
}

// Handlers contains all handler functions
type Handlers struct {
	db         Database
	dispatcher *websockets.Dispatcher
}

// NewHandlers creates and returns a new instance of `Handlers` with the datastore
func NewHandlers(datastore Database, dispatcher *websockets.Dispatcher) *Handlers {
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
	var rawData interface{}
	json.Unmarshal(b, &rawData)
	action := Action{
		Type: "update",
		Path: keys,
		Data: rawData,
	}
	h.dispatcher.Broadcast <- action

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

// WebsockerHandler handles all of the websocket connections
func (h *Handlers) WebsockerHandler(c *gin.Context) {
	serveWs(h.dispatcher, c.Writer, c.Request)
}

// serveWs handles websocket requests from the peer.
func serveWs(dispatcher *websockets.Dispatcher, w http.ResponseWriter, r *http.Request) {
	conn, err := websockets.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := websockets.NewClient(dispatcher, conn)
	client.Dispatcher().Register() <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteDispatch()
	go client.ReadDispatch()
}
