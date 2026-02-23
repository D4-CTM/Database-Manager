package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortErr(c *gin.Context, code int, msg string, err error) {
	log.Printf("[ERROR] %v", err)
	c.AbortWithStatusJSON(code, msg)
}

func abort(c *gin.Context, code int, msg string) {
	log.Printf("[ERROR] %v", msg)
	c.AbortWithStatusJSON(code, msg)
}

func internalError(c *gin.Context, msg string, err error) {
	abortErr(c, http.StatusInternalServerError, msg, err)
}

func unableToBind(c *gin.Context, err error) {
	abortErr(c, http.StatusInternalServerError, "Unable to bind data", err)
}

func dbRefused(c *gin.Context, err error) {
	internalError(c, "Unable to stablish database connection", err)
}

func tablesError(c *gin.Context, err error) {
	internalError(c, "Unable to fetch tables", err)
}
