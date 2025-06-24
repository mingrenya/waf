package handlers // 确保有包声明

import "github.com/gin-gonic/gin"

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "AI-WAF is running"})
}

func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Test endpoint working"})
}

func LoginHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "login success"})
}
