package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()

	server.GET("/hello", func(ctx *gin.Context) {
		a := ctx.ClientIP()
		ctx.String(http.StatusOK, fmt.Sprintf("hello world %s", a))
	})
	server.Run(":8099")
}
