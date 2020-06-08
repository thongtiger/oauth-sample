package store

import (
	"errors"

	"github.com/thongtiger/oauth-rfc6749/auth"
)

// mysqlContext
type mysqlContext struct {
	hostname, port, username, password, database string
}

func NewMYSQLStore(hostname, port, username, password, database string) Store {
	// return interface
	return &mysqlContext{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}

func (c *mysqlContext) NewUser(username, password, role string, scope []string) (*auth.User, error) {
	return nil, errors.New("test")
}
func (c *mysqlContext) ValidateUser(username, password string) (bool, auth.User) {
	return true, auth.User{
		Role:     "admin",
		Scope:    []string{"*"},
		Username: "admin",
	}
}
func (c *mysqlContext) GetUser(username string) (result *auth.User, err error) {
	return nil, errors.New("test error")
}
func (c *mysqlContext) FindUser() (results []*auth.User) {
	return
}
