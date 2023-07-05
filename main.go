package main

import (
	"github.com/gin-gonic/gin"
	logs_server "github.com/razorpay/logsServer/logs-server"
)

func main() {
	r := gin.Default()

	r.GET("/hello", logs_server.HelloController)

	r.Run() // Start the server
}
