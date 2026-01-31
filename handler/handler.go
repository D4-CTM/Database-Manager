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
		w.Write([]byte(err.Error()))
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

	temp.ExecuteTemplate(w, "Connection", map[string]any{
		"Options": []string{
			"Tables",
			"Views",
			"Procedures",
			"Functions",
			"Packages",
			"Sequences",
			"Triggers",
			"Indices",
			"Users",
		},
		"Key": dbName,
	})
}

// Fetches tables owned by database user
func Tables(w http.ResponseWriter, r *http.Request) {
	dbName := r.PathValue("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", dbName))
		return
	}

	query := `
	SELECT
		table_name
	FROM
		sys.all_tables
	WHERE
		owner = :1
	`

	db := cred.GetDB()
	tables := []string{}
	row, err := db.Query(query, strings.ToUpper(cred.User))
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't fetch tables from: %s", dbName))
		return
	}
	defer row.Close()
	t := ""
	for row.Next() {
		row.Scan(&t)
		tables = append(tables, t)
	}

	temp.ExecuteTemplate(w, "Data", map[string]any{
		"data": tables,
		"Opt":  "Table",
		"icon": "table",
	})
}

func Views(w http.ResponseWriter, r *http.Request) {
	dbName := r.PathValue("database")
	cred := service.Cons[dbName]
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't stablish connection with %s", dbName))
		return
	}

	query := `
	SELECT
		view_name
	FROM
		sys.all_views
	WHERE
		owner = :1
	`

	db := cred.GetDB()
	views := []string{}
	row, err := db.Query(query, strings.ToUpper(cred.User))
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		writeStatusMessage(w, http.StatusBadGateway, fmt.Sprintf("Couldn't fetch tables from: %s", dbName))
		return
	}
	defer row.Close()
	v := ""
	for row.Next() {
		row.Scan(&v)
		views = append(views, v)
	}

	temp.ExecuteTemplate(w, "Data", map[string]any{
		"data": views,
		"Opt":  "View",
		"icon": "eye",
	})
}

func Procedures(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Functions(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Packages(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Sequences(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Triggers(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Indices(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Users(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}
