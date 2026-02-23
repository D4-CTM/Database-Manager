package handler

import (
	"dbmt/Service"
	dtos "dbmt/handler/Dtos"
	"fmt"
	"log"
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

	c.JSON(http.StatusOK, cred)
}

func GetConnections(c *gin.Context) {
	cons := make([]string, len(service.Cons))
	idx := 0
	for key := range service.Cons {
		cons[idx] = key
		idx++
	}

	c.JSON(http.StatusOK, cons)
}

func PostConnection(c *gin.Context) {
	cred := service.Credentials{}
	if err := c.ShouldBind(&cred); err != nil {
		internalError(c, "Unable to bind credential", err)
		return
	}

	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
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

	c.JSON(http.StatusOK, conName)
}

func PutConnection(c *gin.Context) {
	oldName := c.Param("oldName")
	conName := c.Query("newName")

	if conName == "" {
		conName = oldName
	}

	if conName != oldName {
		for name := range service.Cons {
			if name == conName && name != oldName {
				abort(c, http.StatusInternalServerError, "Connection name already exists")
				return
			}
		}
	}

	cred := service.Credentials{}
	if err := c.ShouldBind(&cred); err != nil {
		dbRefused(c, err)
		return
	}

	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	if conName != oldName {
		service.Cons[conName] = service.Cons[oldName]
		delete(service.Cons, oldName)
	}
	service.SaveConnections()

	c.JSON(http.StatusOK, gin.H{
		"OldName": oldName,
		"NewName": conName,
	})
}

func DeleteConnection(c *gin.Context) {
	conName := c.Param("conName")
	con := service.Cons[conName]
	if err := con.Close(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	delete(service.Cons, conName)
	service.SaveConnections()
	c.Status(http.StatusNoContent)
}

func Ping(c *gin.Context) {
	dbName := c.Param("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}
	users := []string{}

	data := gin.H{
		"User":     strings.ToUpper(cred.User),
		"Database": cred.Database,
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
		tablesError(c, err)
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
		dbRefused(c, err)
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
	rows, err := db.Query(query, strings.ToUpper(schema))
	if err != nil {
		tablesError(c, err)
		return
	}
	defer rows.Close()

	val := ""
	for rows.Next() {
		rows.Scan(&val)
		values = append(values, val)
	}

	c.JSON(http.StatusOK, values)
}

func Tables(c *gin.Context) {
	fetchData(c, "table_name", "sys.all_tables", "owner = :1")
}

func Views(c *gin.Context) {
	fetchData(c, "view_name", "sys.all_views", "owner = :1")
}

func Procedures(c *gin.Context) {
	fetchData(c, "object_name", "sys.all_procedures", "owner = :1 AND object_type = 'PROCEDURE'")
}

func Functions(c *gin.Context) {
	fetchData(c, "object_name", "sys.all_procedures", "owner = :1 AND object_type = 'FUNCTION'")
}

func Packages(c *gin.Context) {
	fetchData(c, "object_name", "sys.all_procedures", "owner = :1 AND object_type = 'PACKAGE'")
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
	data, err := c.GetRawData()
	if err != nil {
		internalError(c, "Unable to get raw data", err)
		return
	}
	query := string(data)

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	ar, err := cred.Exec(query)
	if err != nil {
		tablesError(c, err)
		return
	}

	c.JSON(http.StatusOK, ar)
}

func Select(c *gin.Context) {
	dbName := c.Param("database")
	data, err := c.GetRawData()
	if err != nil {
		internalError(c, "Unable to get raw data", err)
		return
	}
	query := string(data)

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	table, err := cred.QueryRows(query)
	if err != nil {
		tablesError(c, err)
		return
	}

	c.JSON(http.StatusOK, (*table))
}

func GetTable(c *gin.Context) {
	dbName := c.Param("database")
	gt := dtos.GetTableDto()
	err := c.ShouldBindQuery(&gt)
	if err != nil {
		unableToBind(c, err)
	}

	if err = gt.IsValid(); err != nil {
		internalError(c, err.Error(), err)
		return
	}
	tableName := gt.TableName()

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	table, err := cred.QueryTable(tableName)
	if err != nil {
		tablesError(c, err)
		return
	}

	c.JSON(http.StatusOK, (*table))
}

func GetTableColumnMetadata(c *gin.Context) {
	dbName := c.Param("database")
	gt := dtos.GetTableDto()
	err := c.ShouldBindQuery(&gt)
	if err != nil {
		unableToBind(c, err)
	}

	if err = gt.IsValid(); err != nil {
		internalError(c, err.Error(), err)
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	columns, err := cred.QueryTableColumnsMetadata(gt.Table)
	if err != nil {
		internalError(c, "Unable to fetch columns", err)
		return
	}

	c.JSON(http.StatusOK, columns)
}

func GetTableConstraints(c *gin.Context) {
	dbName := c.Param("database")
	gt := dtos.GetTableDto()
	err := c.ShouldBindQuery(&gt)
	if err != nil {
		unableToBind(c, err)
	}

	if err = gt.IsValid(); err != nil {
		internalError(c, err.Error(), err)
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	columns, err := cred.QueryTableConstraints(gt.Table)
	if err != nil {
		internalError(c, "Unable to fetch columns", err)
		return
	}

	c.JSON(http.StatusOK, columns)
}

func GetFunctionArguments(c *gin.Context) {
	dbName := c.Param("database")
	gf := dtos.GetFunctionArgs()
	err := c.ShouldBindQuery(&gf)
	if err != nil {
		unableToBind(c, err)
	}

	if err = gf.IsValid(); err != nil {
		internalError(c, err.Error(), err)
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	columns, err := cred.QueryFunctionArguments(gf.Owner, gf.ObjectName)
	if err != nil {
		internalError(c, "Unable to fetch function arguments", err)
		return
	}

	c.JSON(http.StatusOK, columns)
}

func GetDdl(c *gin.Context) {
	dbName := c.Param("database")
	gd := dtos.GetDdlDto()
	err := c.ShouldBindQuery(&gd)
	if err != nil {
		unableToBind(c, err)
	}

	if err = gd.IsValid(); err != nil {
		internalError(c, err.Error(), err)
		return
	}

	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	columns, err := cred.QueryDDL(gd.Type, gd.Name, gd.Owner)
	if err != nil {
		internalError(c, "Unable to fetch ddl", err)
		return
	}

	c.JSON(http.StatusOK, columns)
}
