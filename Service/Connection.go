package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
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

func (c *Credentials) connect() error {
	if c.db != nil {
		return nil;
	}

	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
	if err != nil {
		return err;
	}
	c.db = db
	return nil
}

func (c *Credentials) Ping() error {
	if err := c.connect(); err != nil {
		return err
	}

	return c.db.Ping()
}

type Connections map[string]Credentials

func SaveConnections(c Connections, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	j := json.NewEncoder(f)
	j.SetIndent("", "\t")

	if err = j.Encode(c); err != nil {
		return err
	}

	return nil
}

func LoadConnections(path string) (Connections, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := Connections{}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	return c, nil
}
