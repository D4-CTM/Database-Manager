// SlaveMasterPair handler
package handler

import (
	service "dbmt/Service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSlaveMasterPairs(c *gin.Context) {
	pairs := make([]string, len(service.SlaveMaster))
	idx := 0
	for key := range service.Cons {
		pairs[idx] = key
		idx++
	}

	c.JSON(http.StatusOK, pairs)
}

func GetSlaveMasterPair(c *gin.Context) {
	conName := c.Query("conName")
	if conName == "" {
		abort(c, http.StatusBadRequest, "Please provide a valid Connection Name")
		return
	}

	pair := service.SlaveMaster[conName]

	c.JSON(http.StatusOK, pair)
}

func DeleteSlaveMasterPair(c *gin.Context) {
	conName := c.Param("conName")
	pair := service.SlaveMaster[conName]
	if err := pair.Close(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	delete(service.SlaveMaster, conName)
	service.SaveSlaveMasterPair()
	c.Status(http.StatusNoContent)
}

func PostSlaveMasterPair(c *gin.Context) {
	pair := service.SlaveMasterPair{}
	if err := c.ShouldBind(&pair); err != nil {
		internalError(c, "Unable to bind credential", err)
		return
	}
	
	if err := pair.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	dbName := fmt.Sprintf("%s_%s", pair.Master.Database, pair.Slave.Database)
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

	service.SlaveMaster[conName] = pair
	service.SaveSlaveMasterPair()

	c.JSON(http.StatusOK, conName)
}

func PutSlaveMasterPair(c *gin.Context) {
	oldName := c.Param("oldName")
	conName := c.Query("newName")

	if conName == "" {
		abort(c, http.StatusInternalServerError, "New connection name cannot be empty!",)
		return 
	}

	if conName != oldName {
		for name := range service.Cons {
			if name == conName && name != oldName {
				abort(c, http.StatusInternalServerError, "Connection name already exists")
				return
			}
		}
	}

	pair := service.SlaveMasterPair{}
	if err := c.ShouldBind(&pair); err != nil {
		internalError(c, "Unable to bind credential", err)
		return
	}

	if err := pair.Ping(); err != nil {
		dbRefused(c, err)
		return
	}

	if conName != oldName {
		service.SlaveMaster[conName] = service.SlaveMaster[oldName]
		delete(service.Cons, oldName)
	}
	service.SaveSlaveMasterPair()

	c.JSON(http.StatusOK, gin.H{
		"OldName": oldName,
		"NewName": conName,
	})
}
