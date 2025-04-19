package main

import (
	"github.com/gin-gonic/gin"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
	cronv3 "github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []saramax.Consumer
	cron      *cronv3.Cron
}
