package handler

import (
	"dbmt/Service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var temp = template.Must(template.ParseFiles("templates/base.html", "templates/ping.html"))

func writeStatusMessage(w http.ResponseWriter, status int, message string) {
	log.Printf("[INFO] Writing status message: %s\n", message)
	w.Header().Set("HX-Message", message)
	w.WriteHeader(status)
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Connections": service.Cons,
	}

	if err := temp.Execute(w, data); err != nil {
		log.Printf("%v", err)
		w.Write([]byte(err.Error()))
		writeStatusMessage(w, http.StatusBadGateway, "Error executing the template")
	}
}

func CreateConnection(w http.ResponseWriter, r *http.Request) {
	i, _ := strconv.Atoi(r.PostFormValue("Port"))
	cred := service.Credentials{
		Server:   r.PostFormValue("Server"),
		Port:     i,
		Database: r.PostFormValue("Database"),
		User:     r.PostFormValue("Username"),
		Password: r.PostFormValue("Password"),
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

	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR on Connect()] %v", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", conName))
		return
	}

	service.Cons[conName] = cred

	temp.ExecuteTemplate(w, "ConnectionList", map[string]any{
		"Connections": service.Cons,
	})
}

func Ping(w http.ResponseWriter, r *http.Request) {
	dbName := r.PathValue("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", dbName))
		return
	}

	query := `
	SELECT
		username
	FROM
		sys.all_users;
	`

	db := cred.GetDB()
	users := []string{}
	row, err := db.Query(query)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't fetch users from: %s", dbName))
		return
	}
	defer row.Close()
	u := ""
	for row.Next() {
		row.Scan(&u)
		if (u == cred.User) {
			continue
		}
		users = append(users, u)
	}

	temp.ExecuteTemplate(w, "Schemas", map[string]any{
		"Schemas": users,
		"User": strings.ToUpper(cred.User),
		"Database": cred.Database,
		"Key": dbName,
	})
}

func Options(w http.ResponseWriter, r *http.Request) {
	dbName := r.PathValue("database")
	schema := r.PathValue("schema")

	temp.ExecuteTemplate(w, "Options", map[string]any{
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
		"Key": dbName,
	})
}

// Params:
// ss - Select Statement
// table - Db Table
// wc - where condition
func fetchData(w http.ResponseWriter, r *http.Request, ss string, table string, wc string, data map[string]any) {
	dbName := r.PathValue("database")
	schema := r.PathValue("schema")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", dbName))
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
	trigger := []string{}
	row, err := db.Query(query, strings.ToUpper(schema))
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't fetch %s from: %s", ss, dbName))
		return
	}
	defer row.Close()
	t := ""
	for row.Next() {
		row.Scan(&t)
		trigger = append(trigger, t)
	}

	data["data"] = trigger
	temp.ExecuteTemplate(w, "Data", data)
}

// Fetches tables owned by database user
func Tables(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Table",
		"icon": "table",
	}
	fetchData(w, r, "table_name", "sys.all_tables", "owner = :1", data)
}

func Views(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "View",
		"icon": "eye",
	}
	fetchData(w, r, "view_name", "sys.all_views", "owner = :1", data)
}

func Procedures(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Procedures",
		"icon": "code",
	}
	fetchData(w, r, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PROCEDURE'", data)
}

func Functions(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Functions",
		"icon": "code",
	}
	fetchData(w, r, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'FUNCTION'", data)
}

func Packages(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Packages",
		"icon": "linode",
	}
	fetchData(w, r, "procedure_name", "sys.all_procedures", "owner = :1 AND object_type = 'PACKAGE'", data)
}

func Sequences(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Sequences",
		"icon": "line-chart",
	}
	fetchData(w, r, "sequence_name", "sys.all_sequences", "sequence_owner = :1", data)
}

func Triggers(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Triggers",
		"icon": "exchange",
	}
	fetchData(w, r, "trigger_name", "sys.all_triggers", "owner = :1", data)
}

func Indexes(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Opt":  "Triggers",
		"icon": "exchange",
	}
	fetchData(w, r, "index_name", "sys.all_indexes", "owner = :1", data)
}

func Users(w http.ResponseWriter, r *http.Request) {
	dbName := r.PathValue("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", dbName))
		return
	}

	query := `
	SELECT
		username
	FROM
		sys.all_users;
	`

	db := cred.GetDB()
	users := []string{}
	row, err := db.Query(query)
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't fetch users from: %s", dbName))
		return
	}
	defer row.Close()
	u := ""
	for row.Next() {
		row.Scan(&u)
		users = append(users, u)
	}

	temp.ExecuteTemplate(w, "Data", map[string]any{
		"data": users,
		"Opt":  "Users",
		"icon": "user",
	})
}
