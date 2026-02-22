package service

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// Ui indicator
	ShowAll bool

	db      *sql.DB
}

type Table struct {
	Name        string
	ColumnNames []string
	Rows        [][]any
}

func CreateCredFromGin(c *gin.Context) Credentials {
	i, _ := strconv.Atoi(c.Request.PostFormValue("Port"))
	showAll := c.Request.PostFormValue("ShowAll")
	return Credentials{
		Server:   c.Request.PostFormValue("Server"),
		Port:     i,
		Database: c.Request.PostFormValue("Database"),
		User:     c.Request.PostFormValue("Username"),
		Password: c.Request.PostFormValue("Password"),
		ShowAll:  showAll == "on",
	}

}

func (c *Credentials) Connect() error {
	if c.db != nil {
		return nil
	}

	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
	if err != nil {
		return fmt.Errorf("Unable to stablish connection!\n%v", err)
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

func (c *Credentials) Exec(query string) (int64, error) {
	db := c.GetDB()
	result, err := db.Exec(query)
	if err != nil {
		return -1, fmt.Errorf("[ERROR] %v\n", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("[ERROR] %v\n", err)
	}

	return affected, nil
}

func (c *Credentials) QueryRows(query string) (*Table, error) {
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v\n", err)
	}
	defer rows.Close()

	var table = Table{}
	err = fetchRows(rows, &table)

	return &table, err
}

func (c *Credentials) QueryTable(tableName string) (*Table, error) {
	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] %v\n", err)
	}
	defer rows.Close()

	var table = Table{
		Name:        tableName,
	}
	err = fetchRows(rows, &table)

	return &table, err
}

func fetchRows(rows *sql.Rows, table *Table) (error) {
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[WARN] Column couldn't be read!\n\t[ERROR] Error %v", err)
	}
	colLen := len(cols)

	table.ColumnNames = cols
	table.Rows = [][]any{{}}

	for rows.Next() {
		columns := make([]any, colLen)
		columnPointers := make([]any, colLen)
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return fmt.Errorf("[ERROR] %v", err)
		}

		rowData := make([]any, colLen)
		for i := range cols {
			val := columns[i]

			if b, ok := val.([]byte); ok {
				rowData[i] = string(b)
			} else {
				rowData[i] = val
			}
			fmt.Printf("%v ", rowData[i])
		}
		table.Rows = append(table.Rows, rowData)
		fmt.Println()
	}

	return nil
}
