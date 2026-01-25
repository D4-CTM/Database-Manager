package handler

import (
	"log"
	"net/http"
	"text/template"
)

var temp = template.Must(template.ParseFiles("index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	if err := temp.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		w.Write([]byte(err.Error()))
	}
}
