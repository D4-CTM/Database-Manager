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

func Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Connections": service.Cons,
	}
	if err := temp.Execute(w, data); err != nil {
		log.Printf("%v", err)
		w.Write([]byte(err.Error()))
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
				conName = fmt.Sprintf("%s %d", dbName, idx)
				idx++
			}
		}
	}

	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR on Connect()] %v", err)
		w.Write([]byte(err.Error()))
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
	var errStr string = ""
	if err := cred.Ping(); err != nil {
		log.Printf("[ERROR] %v\n", err)
		errStr = err.Error()
	}

	temp.ExecuteTemplate(w, "Connection", map[string]any{
		"ConnectionStablished": errStr == "",
		"Error": errStr,
		"Key": dbName,
	})
}
