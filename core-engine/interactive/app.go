package main

import (
	"TestCopilot/TestEngine/pkg/grpcx"
	"TestCopilot/TestEngine/pkg/saramax"
)

type App struct {
	// 所有需要 main 函数启动，关闭的，都会在这里有一个
	server    *grpcx.Server
	consumers []saramax.Consumer
}
