package startup

import (
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
