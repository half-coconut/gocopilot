package ioc

import (
	"TestCopilot/backend/pkg/logger"
	"go.uber.org/zap"
)

// InitLogger 这是全局的log，使用接口的依赖注入，将 logger.LoggerV1 注入log
// 方便后期可以扩展或者更换 zap
func InitLogger() logger.LoggerV1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
