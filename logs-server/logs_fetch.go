package logs_server

import (
	"github.com/gin-gonic/gin"
)

func HelloController(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello, World!"})
}
