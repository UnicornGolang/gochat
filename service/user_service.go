package service

import (
	"gochat/models"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
  data := models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

func GetUser(c *gin.Context) {
  user := models.GetUser()
	c.JSON(200, gin.H{
		"message": user,
	})
}
