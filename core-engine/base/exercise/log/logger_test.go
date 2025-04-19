package log

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

func TestVipper(t *testing.T) {
	logPath()
}

func TestSugaredLogger(t *testing.T) {
	InitSugaredLogger()
	err := errors.New("这是错误信息")
	SL.Info("这是一条测试数据", zap.Error(err), zap.String("name", "hello base"))
}

func TestLogger(t *testing.T) {
	InitLogger()
	err := errors.New("L: 这是错误信息")
	L.Info("L: 这是一条测试数据", zap.Error(err), zap.String("name", "hello base"))
}
