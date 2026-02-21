package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code     int
	ErrorMsg string
}

func abort(c *gin.Context, code int, msg string) {
	log.Printf("[ERROR] %v", msg)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"Error": Error{
			Code: code,
			ErrorMsg: msg,
		},
	})
}

func abortWithError(c *gin.Context, err error, code int, msg string) {
	log.Printf("[ERROR] %v", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"Error": Error{
			Code: code,
			ErrorMsg: msg,
		},
	})
}

func internalError(c *gin.Context, err error, msg string) {
	abortWithError(c, err, http.StatusInternalServerError, msg)
}

func dbConnectionRefused(c *gin.Context, err error) {
	abortWithError(c, err, http.StatusInternalServerError, "Unable to connect with database")
}
