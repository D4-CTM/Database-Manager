package main

import (
	"context"
	"dbmt/Service"
	"dbmt/handler"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ADDR string = ":5461"
	JSON_PATH string = "data.json"
)

func main() {
	err := service.LoadConnections(JSON_PATH)
	if err != nil {
		log.Printf("%v", err)
	}
	defer service.SaveConnections(JSON_PATH)

	log.Println("Starting server...")
	static := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", static))
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/Create/", handler.CreateConnection)
	http.HandleFunc("/Ping/{database}", handler.Ping)
	http.HandleFunc("/Tables/{database}", handler.Tables)
	http.HandleFunc("/Views/{database}", handler.Views)
	http.HandleFunc("/Procedures/{database}", handler.Procedures)
	http.HandleFunc("/Functions/{database}", handler.Functions)
	http.HandleFunc("/Packages/{database}", handler.Packages)
	http.HandleFunc("/Sequences/{database}", handler.Sequences)
	http.HandleFunc("/Triggers/{database}", handler.Triggers)
	http.HandleFunc("/Indices/{database}", handler.Indices)
	http.HandleFunc("/Users/{database}", handler.Users)

	server := &http.Server{
		Addr: ADDR,
	}

    go func() {
		log.Printf("Server running at: http://localhost%s\n", ADDR)
        if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownRelease()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("HTTP shutdown error: %v", err)
    }
}
