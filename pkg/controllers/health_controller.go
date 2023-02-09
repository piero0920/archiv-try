package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piero0920/archiv-try/pkg/database"
)

func GetHealth(c *gin.Context) {
	db, _ := database.DB.DB()

	if err := db.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "no connection to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "Ok",
	})
}
