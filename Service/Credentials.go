package service

import (
	"database/sql"
	"fmt"

	_ "github.com/godror/godror"
)

// Oracle credentials
type Credentials struct {
	// Connection data
	Database string
	Server   string
	Port     int

	// Credentials
	Password string
	User     string

	db *sql.DB
}

func (c *Credentials) Connect() error {
	if c.db != nil {
		return nil;
	}

	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
	if err != nil {
		return fmt.Errorf("Unable to stablish connection!\n%v", err);
	}
	c.db = db
	return nil
}

func (c *Credentials) Ping() error {
	if c.db == nil {
		if err := c.Connect(); err != nil {
			return nil
		}
	}

	return c.db.Ping()
}

func (c *Credentials) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *Credentials) GetDB() *sql.DB {
	return c.db
}
