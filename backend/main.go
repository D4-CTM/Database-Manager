package main

import (
	"context"
	"dbmt/Service"
	"dbmt/handler"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADDR string = ":5461"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"0.0.0.0"})

	JSON_PATH := path.Join(os.Getenv("CREDS_SUBDIR"), "data.json")
	err := service.LoadConnections(JSON_PATH)
	if err != nil {
		log.Printf("%v", err)
	}
	defer service.SaveConnections(JSON_PATH)

	log.Println("Starting server...")
	r.GET("/", handler.Index)
	r.POST("/Create/", handler.CreateConnection)
	r.GET("/Ping/{database}/", handler.Ping)
	r.GET("/Options/{database}/{schema}", handler.Options)

	r.GET("/Tables/{database}/{schema}", handler.Tables)
	r.GET("/Views/{database}/{schema}", handler.Views)
	r.GET("/Procedures/{database}/{schema}", handler.Procedures)
	r.GET("/Functions/{database}/{schema}", handler.Functions)
	r.GET("/Packages/{database}/{schema}", handler.Packages)
	r.GET("/Sequences/{database}/{schema}", handler.Sequences)
	r.GET("/Triggers/{database}/{schema}", handler.Triggers)
	r.GET("/Indices/{database}/{schema}", handler.Indexes)

	r.GET("/Query/{database}", handler.Query)

    go func() {
	  log.Printf("Server running at: http://localhost%s\n", ADDR)
      if err := r.Run(ADDR); err != nil {
          log.Fatalf("HTTP server error: %v", err)
      }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    _, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownRelease()	
}
