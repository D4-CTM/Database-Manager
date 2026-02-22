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
	r.GET("/api/Connection/list", handler.GetConnections)
	r.GET("/api/Connection", handler.GetCredential)
	r.POST("/api/Connection", handler.PostConnection)
	r.PUT("/api/Connection/:oldName", handler.PutConnection)

	r.GET("/api/Ping/:database", handler.Ping)
	r.GET("/api/Tables/:database/:schema", handler.Tables)
	r.GET("/api/Views/:database/:schema", handler.Views)
	r.GET("/api/Procedures/:database/:schema", handler.Procedures)
	r.GET("/api/Functions/:database/:schema", handler.Functions)
	r.GET("/api/Packages/:database/:schema", handler.Packages)
	r.GET("/api/Sequences/:database/:schema", handler.Sequences)
	r.GET("/api/Triggers/:database/:schema", handler.Triggers)
	r.GET("/api/Indices/:database/:schema", handler.Indexes)

	r.GET("/api/Select/:database", handler.Select)
	r.GET("/api/Exec/:database", handler.Exec)
	r.GET("/api/Query/:database", handler.GetTable)

	log.Printf("Server running at: http://localhost%s\n", ADDR)
	if err := r.Run(ADDR); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
