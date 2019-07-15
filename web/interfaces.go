package web

// Database is the required interface to the DB layer from the HTTP handlers
type Database interface {
	CreateProject(projectName string, data []byte) error
	GetProject(projectName string) ([]byte, error)
	DeleteProject(projectName string) error

	GetProjectKey(projectName string, keys ...string) ([]byte, error)
	CreateProjectKey(projectName string, data []byte, keys ...string) error
	UpdateProjectKey(projectName string, data []byte, keys ...string) error
	DeleteProjectKey(projectName string, keys ...string) error
}
