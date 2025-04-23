package web

import "github.com/gin-gonic/gin"

type Result struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type handler interface {
	RegisterRoutes(server *gin.Engine)
}
