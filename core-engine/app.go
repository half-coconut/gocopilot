package main

import (
	"TestCopilot/TestEngine/pkg/saramax"
	"github.com/gin-gonic/gin"
	cronv3 "github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []saramax.Consumer
	cron      *cronv3.Cron
}
