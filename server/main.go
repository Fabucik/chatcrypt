package main

import (
	"github.com/Fabucik/chatcrypt/server/messaging"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.POST("/sendmessage", messaging.SendMessage)
}
