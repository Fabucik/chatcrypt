package messaging

import (
	"github.com/Fabucik/chatcrypt/entities"
	"github.com/gin-gonic/gin"
)

func SendMessage(ctx *gin.Context) {
	var senderInfo entities.Sender
	ctx.Bind(&senderInfo)
}
