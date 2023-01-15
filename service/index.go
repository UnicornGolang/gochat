package service

import "github.com/gin-gonic/gin"

// @BasePath /api/v1
// @Accept json
// @Produce json
// PingExample godoc
// @Summary ping example
// @Schemes

// @Description 冒烟测试
// @Tags 首页
// @Success 200 {string} welcome.
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "welcome!!",
	})
}
