package main

import (
	"TestCopilot/TestEngine/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server    *gin.Engine
	Consumers []events.Consumer
}
