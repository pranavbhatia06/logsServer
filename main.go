package main

import (
	"github.com/gin-gonic/gin"
	logs_server "github.com/razorpay/logsServer/logs-server"
)

func main() {
	r := gin.Default()

	r.POST("/logs", logs_server.GetLogs)

	r.Run() // Start the server
}
