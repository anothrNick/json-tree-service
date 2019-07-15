package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Postgres wraps the db connection to a postgres instance
type Postgres struct {
	db *sql.DB
}

// NewPostgres returns a new instance of `Postgres` with a connection to the postgres db based on provided parameters.
func NewPostgres(user, password, host, database string) (*Postgres, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db: db,
	}, nil
}

// TranslateError translates the database specific error to a simple `error` to return to the user.
func (p *Postgres) TranslateError(err error) error {
	fmt.Println(err)
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		default:
			return errors.New(strings.TrimPrefix(err.Error(), "pq: "))
		}
	}

	return errors.New("unknown error")
}

// CreateProject creates a new project with a JSON tree
func (p *Postgres) CreateProject(projectName string, data []byte) error {
	_, err := p.db.Exec("INSERT INTO trees (project, data) VALUES ($1, $2)", projectName, data)
	return err
}

//DeleteProject permanently deletes an entire project's tree
func (p *Postgres) DeleteProject(projectName string) error {
	_, err := p.db.Exec("DELETE FROM trees where project=$1", projectName)
	return err
}

// GetProjectKey retrieves the object at the key path
func (p *Postgres) GetProjectKey(projectName string, keys ...string) ([]byte, error) {
	byt := []byte{}

	keysFormat := strings.Join(keys, ",")
	err := p.db.QueryRow(
		fmt.Sprintf("SELECT data#>'{%s}' as data FROM trees WHERE project=$1 ORDER BY id DESC LIMIT 1", keysFormat),
		projectName).Scan(&byt)

	return byt, err
}

// CreateProjectKey saves the data at the provided key path. Fails if the key already exists.
func (p *Postgres) CreateProjectKey(projectName string, data []byte, keys ...string) error {
	keysFormat := strings.Join(keys, ",")
	_, err := p.db.Exec(fmt.Sprintf("UPDATE trees set data=jsonb_insert(data, '{%s}', $1) WHERE project=$2", keysFormat), data, projectName)
	return err
}

// UpdateProjectKey updates the data at the key path. Creates a new key if it does not already exist.
func (p *Postgres) UpdateProjectKey(projectName string, data []byte, keys ...string) error {
	keysFormat := strings.Join(keys, ",")
	_, err := p.db.Exec(fmt.Sprintf("UPDATE trees set data=jsonb_set(data, '{%s}', $1) WHERE project=$2", keysFormat), data, projectName)
	return err
}

// DeleteProjectKey permanently removes the data at the key path.
func (p *Postgres) DeleteProjectKey(projectName string, keys ...string) error {
	keysFormat := strings.Join(keys, ",")
	_, err := p.db.Exec(fmt.Sprintf("UPDATE trees SET data=data #- '{%s}' where project=$1", keysFormat), projectName)
	return err
}
