package main

import (
	"dbmt/Service"
	"dbmt/handler"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	ADDR string = ":5461"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"0.0.0.0"})

	err := service.LoadConnections()
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Starting server...")
	r.GET("/Connection", handler.GetConnections)
	r.POST("/Connection", handler.PostConnection)
	r.PUT("/Connection/:oldName", handler.PutConnection)

	r.GET("/Ping/:database", handler.Ping)
	r.GET("/Tables/:database/:schema", handler.Tables)
	r.GET("/Views/:database/:schema", handler.Views)
	r.GET("/Procedures/:database/:schema", handler.Procedures)
	r.GET("/Functions/:database/:schema", handler.Functions)
	r.GET("/Packages/:database/:schema", handler.Packages)
	r.GET("/Sequences/:database/:schema", handler.Sequences)
	r.GET("/Triggers/:database/:schema", handler.Triggers)
	r.GET("/Indices/:database/:schema", handler.Indexes)

	r.GET("/Query/:database", handler.GetTable)

	log.Printf("Server running at: http://localhost%s\n", ADDR)
	if err := r.Run(ADDR); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
