package startup

import (
	"TestCopilot/TestEngine/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
