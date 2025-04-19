package main

import (
	"github.com/half-coconut/gocopilot/core-engine/pkg/grpcx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/saramax"
)

type App struct {
	// 所有需要 main 函数启动，关闭的，都会在这里有一个
	server    *grpcx.Server
	consumers []saramax.Consumer
}
