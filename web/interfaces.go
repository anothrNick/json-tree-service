package web

import "github.com/anothrNick/json-tree-service/database"

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
