package mylog_t

import (
	"TestCopilot/TestEngine/base/example/mylog"
	"context"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	mylog.Init(nil)
	defer mylog.Sync() // Sync 将缓存中的日志刷新到磁盘文件中

	c, cancle := context.WithTimeout(context.Background(), time.Minute)
	defer cancle()

	mylog.C(c).Warnw("Create post function called")
}

func logOptions() *mylog.Options {
	return &mylog.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
