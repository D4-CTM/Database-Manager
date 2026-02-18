package handler

import (
	"dbmt/Service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResultValue int

const (
	TableData ResultValue = iota
)

func writeStatusMessage(w http.ResponseWriter, status int, message string) {
	log.Printf("[INFO] Writing status message: %s\n", message)
	w.Header().Set("HX-Message", message)
	w.WriteHeader(status)
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Connections": service.Cons,
	})
}

func CreateConnection(c *gin.Context) {
	i, _ := strconv.Atoi(c.PostForm("Port"))
	showAll := c.PostForm("ShowAll")
	cred := service.Credentials{
		Server:   c.DefaultPostForm("Server", "localhost"),
		Port:     i,
		Database: c.PostForm("Database"),
		User:     c.PostForm("Username"),
		Password: c.PostForm("Password"),
		ShowAll:  showAll == "on",
	}

	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR on Connect()] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dbName := cred.Database
	conName := dbName
	{
		idx := 1
		for name := range service.Cons {
			if name == conName {
				conName = fmt.Sprintf("%s_%d", dbName, idx)
				idx++
			}
		}
	}

	service.Cons[conName] = cred

	c.JSON(http.StatusOK, gin.H{
		"Connections": service.Cons,
	})
}

func Ping(c *gin.Context) {
	dbName := c.Request.PathValue("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR on Ping()] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	users := []string{}

	data := gin.H{
		"User":     strings.ToUpper(cred.User),
		"Database": cred.Database,
		"Key":      dbName,
	}
	if !cred.ShowAll {
		users = append(users, cred.User)

		data["Schemas"] = users
		c.JSON(http.StatusOK, data)
		return
	}

	query := `
	SELECT
		username
	FROM
		sys.all_users;
	`

	db := cred.GetDB()
	row, err := db.Query(query)
	if err != nil {
		log.Printf("[ERROR on Ping()] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer row.Close()
	u := ""
	for row.Next() {
		row.Scan(&u)
		if u == cred.User {
			continue
		}
		users = append(users, u)
	}

	data["Schemas"] = users
	c.JSON(http.StatusOK, data)
}

func Options(c *gin.Context) {
	dbName := c.Request.PathValue("database")
	schema := c.Request.PathValue("schema")

	
	c.JSON(http.StatusOK, gin.H{
		"Options": []string{
			"Tables",
			"Views",
			"Procedures",
			"Functions",
			"Packages",
			"Sequences",
			"Triggers",
			"Indices",
		},
		"Schema": schema,
		"Key":    dbName,
	})
}

// Params:
// ss - Select Statement
// table - Db Table
// wc - where condition
func fetchData(c *gin.Context, ss string, table string, wc string, data map[string]any) {
	dbName := c.Request.PathValue("database")
	schema := c.Request.PathValue("schema")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	query := fmt.Sprintf(`
	SELECT
		%s
	FROM
		%s
	WHERE
		%s;
	`, ss, table, wc)

	db := cred.GetDB()
	values := []string{}
	row, err := db.Query(query, strings.ToUpper(schema))
	if err != nil {
		log.Printf("[ERROR] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer row.Close()
	val := ""
	for row.Next() {
		row.Scan(&val)
		values = append(values, val)
	}

	data["data"] = values
	data["Schema"] = schema
	data["ConName"] = dbName
	c.JSON(http.StatusOK, data)
}

func Tables(c *gin.Context) {
	data := map[string]any{
		"icon": "table",
	}
	fetchData(c, "table_name", "sys.all_tables", "owner = :1", data)
}

func Views(c *gin.Context) {
	data := map[string]any{
		"icon": "eye",
	}
	fetchData(c, "view_name", "sys.all_views", "owner = :1", data)
}

func Procedures(c *gin.Context) {
	data := map[string]any{
		"icon": "code",
	}
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PROCEDURE'", data)
}

func Functions(c *gin.Context) {
	data := map[string]any{
		"icon": "code",
	}
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'FUNCTION'", data)
}

func Packages(c *gin.Context) {
	data := map[string]any{
		"icon": "linode",
	}
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PACKAGE'", data)
}

func Sequences(c *gin.Context) {
	data := map[string]any{
		"icon": "line-chart",
	}
	fetchData(c, "sequence_name", "sys.all_sequences", "sequence_owner = :1", data)
}

func Triggers(c *gin.Context) {
	data := map[string]any{
		"icon": "exchange",
	}
	fetchData(c, "trigger_name", "sys.all_triggers", "owner = :1", data)
}

func Indexes(c *gin.Context) {
	data := map[string]any{
		"icon": "exchange",
	}
	fetchData(c, "index_name", "sys.all_indexes", "owner = :1", data)
}

func Query(c *gin.Context) {
	dbName := c.Request.PathValue("database")
	urlQuery := c.Request.URL.Query()
	tableName := urlQuery.Get("Table")
	owner := strings.TrimSpace(urlQuery.Get("Owner"))
	if len(owner) > 0 {
		tableName = fmt.Sprintf("%s.%s", owner, tableName)
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	table, err := cred.QueryTable(tableName)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ResultData": TableData,
		"Table":      (*table),
	})
}
