package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var SL *zap.SugaredLogger
var L *zap.Logger

func InitSugaredLogger() {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath(),
		MaxSize:    100,
		MaxBackups: 5,
		Compress:   false,
	}
	consoleSyncer := zapcore.AddSync(os.Stdout)
	writeSyncer := zapcore.AddSync(lumberJackLogger)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleSyncer, zap.DebugLevel),
		zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel),
	)
	SL = zap.New(core, zap.AddCaller()).Sugar()
}

func InitLogger() {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath(),
		MaxSize:    100,
		MaxBackups: 5,
		Compress:   false,
	}
	consoleSyncer := zapcore.AddSync(os.Stdout)
	writeSyncer := zapcore.AddSync(lumberJackLogger)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleSyncer, zap.DebugLevel),
		zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel),
	)
	L = zap.New(core, zap.AddCaller())

	defer L.Sync()
}

func logPath() string {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetString("log.path")
}
