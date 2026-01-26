package handler

import (
	"dbmt/Service"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func Tables(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
}

func Views(w http.ResponseWriter, r *http.Request) {
	writeStatusMessage(w, http.StatusNotImplemented, "Not implemented")
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
