package main

import (
	"context"
	"dbmt/Service"
	"dbmt/handler"
	"log"
	"os"
	"os/signal"
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

	err := service.LoadConnections()
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("Starting server...")
	r.GET("/api/Connection/list", handler.GetConnections)
	r.GET("/api/Connection", handler.GetCredential)
	r.POST("/api/Connection", handler.PostConnection)
	r.PUT("/api/Connection/:oldName", handler.PutConnection)
	r.DELETE("/api/Connection/:conName", handler.DeleteConnection)

	r.GET("/api/Ping/:database", handler.Ping)
	r.GET("/api/Tables/:database/:schema", handler.Tables)
	r.GET("/api/Views/:database/:schema", handler.Views)
	r.GET("/api/Procedures/:database/:schema", handler.Procedures)
	r.GET("/api/Functions/:database/:schema", handler.Functions)
	r.GET("/api/Packages/:database/:schema", handler.Packages)
	r.GET("/api/Sequences/:database/:schema", handler.Sequences)
	r.GET("/api/Triggers/:database/:schema", handler.Triggers)
	r.GET("/api/Indices/:database/:schema", handler.Indexes)

	r.POST("/api/Select/:database", handler.Select)
	r.PUT("/api/Exec/:database", handler.Exec)
	r.GET("/api/Query/Table/:database", handler.GetTable)
	r.GET("/api/Query/Columns/:database", handler.GetTableColumnMetadata)
	r.GET("/api/Query/Constraints/:database", handler.GetTableConstraints)
	r.GET("/api/Query/Function/Arguments/:database", handler.GetFunctionArguments)
	r.GET("/api/Query/Package/Body/:database", handler.GetPackageBody)
	r.GET("/api/Query/DDL/:database", handler.GetDdl)

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
	service.Cons.Close()
}
