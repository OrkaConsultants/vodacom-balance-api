package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HelloController struct{}

func (u HelloController) SayHello(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Hello from " + viper.GetString("app.name") + " !!!!"})
	ctx.Abort()
	return
}
