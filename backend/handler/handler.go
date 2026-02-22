package handler

import (
	"dbmt/Service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCredential(c *gin.Context) {
	conName := c.Query("conName")
	if conName == "" {
		abort(c, http.StatusBadRequest, "Please provide a valid Connection Name")
		return
	}

	cred := service.Cons[conName]

	c.JSON(http.StatusOK, gin.H{
		"Credential": cred,
	})
}

func GetConnections(c *gin.Context) {
	cons := make([]string, len(service.Cons))
	idx := 0
	for key := range service.Cons {
		cons[idx] = key
		idx++
	}

	c.JSON(http.StatusOK, gin.H{
		"Connections": cons,
	})
}

func PostConnection(c *gin.Context) {
	cred := service.Credentials{}
	if err := c.ShouldBind(&cred); err != nil {
		internalError(c, err, "Unable to bind credential")
		return
	}

	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
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
	service.SaveConnections()

	c.JSON(http.StatusOK, gin.H{
		"Name": conName,
	})
}

func PutConnection(c *gin.Context) {
	oldName := c.Param("oldName")
	conName := c.Query("newName")

	if conName == "" {
		conName = oldName
	}

	if conName != oldName {
		for name := range service.Cons {
			if name == conName {
				abort(c, http.StatusInternalServerError, "Connection name already exists")
				return
			}
		}
	}

	cred := service.Credentials{}
	if err := c.ShouldBind(&cred); err != nil {
		dbConnectionRefused(c, err)
		return
	}

	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
		return
	}

	if conName != oldName {
		service.Cons[conName] = service.Cons[oldName]
		delete(service.Cons, oldName)
	}
	service.SaveConnections()

	c.JSON(http.StatusOK, gin.H{
		"Name": conName,
	})
}

func Ping(c *gin.Context) {
	dbName := c.Param("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
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
		internalError(c, err, "Unable to fetch user table")
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

// Params:
// ss - Select Statement
// table - Db Table
// wc - where condition
func fetchData(c *gin.Context, ss string, table string, wc string) {
	dbName := c.Param("database")
	schema := c.Param("schema")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
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
		internalError(c, err, "Unable to fetch tables")
		return
	}
	defer row.Close()
	val := ""
	for row.Next() {
		row.Scan(&val)
		values = append(values, val)
	}

	data := gin.H{}
	data["data"] = values
	data["Schema"] = schema
	data["ConName"] = dbName
	c.JSON(http.StatusOK, data)
}

func Tables(c *gin.Context) {
	fetchData(c, "table_name", "sys.all_tables", "owner = :1")
}

func Views(c *gin.Context) {
	fetchData(c, "view_name", "sys.all_views", "owner = :1")
}

func Procedures(c *gin.Context) {
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PROCEDURE'")
}

func Functions(c *gin.Context) {
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'FUNCTION'")
}

func Packages(c *gin.Context) {
	fetchData(c, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PACKAGE'")
}

func Sequences(c *gin.Context) {
	fetchData(c, "sequence_name", "sys.all_sequences", "sequence_owner = :1")
}

func Triggers(c *gin.Context) {
	fetchData(c, "trigger_name", "sys.all_triggers", "owner = :1")
}

func Indexes(c *gin.Context) {
	fetchData(c, "index_name", "sys.all_indexes", "owner = :1")
}


func Exec(c *gin.Context) {
	dbName := c.Param("database")
	query := c.Query("query")
	if query == "" {
		abort(c, http.StatusInternalServerError, "Query not specified")
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
		return
	}

	ar, err := cred.Exec(query)
	if err != nil {
		internalError(c, err, "Unable to fetch tables")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AffectedRows": ar,
	})
}

func Select(c *gin.Context) {
	dbName := c.Param("database")
	query := c.Query("query")
	if query == "" {
		abort(c, http.StatusInternalServerError, "Query not specified")
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
		return
	}

	table, err := cred.QueryRows(query)
	if err != nil {
		internalError(c, err, "Unable to fetch tables")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AffectedRows": (*table),
	})
}

func GetTable(c *gin.Context) {
	dbName := c.Param("database")
	tableName := c.Query("Table")
	if tableName == "" {
		abort(c, http.StatusInternalServerError, "Table not specified")
		return
	}

	owner := strings.TrimSpace(c.Query("Owner"))
	if len(owner) > 0 {
		tableName = fmt.Sprintf("%s.%s", owner, tableName)
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbConnectionRefused(c, err)
		return
	}

	table, err := cred.QueryTable(tableName)
	if err != nil {
		internalError(c, err, "Unable to fetch tables")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Table": (*table),
	})
}
