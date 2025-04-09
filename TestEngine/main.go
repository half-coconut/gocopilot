package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"net/http"
)

func main() {
	initViperV1()
	initLogger()
	initPrometheus()

	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	// 定时任务的开启
	//app.cron.Start()

	server := app.server
	server.Run(":3002")

	//// 一分钟内要关完，要退出
	//ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()
	//log.Println(ctx)
	//// 等待运行完毕
	//ctx = app.cron.Stop()
	//// 超时强制退出，防止有些任务执行时间过长
	//tm := time.NewTimer(time.Minute * 10)
	//select {
	//case <-tm.C:
	//case <-ctx.Done():
	//}

}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace，直接用 zap.L()，你啥都打不出来。
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")

	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}

func initViperV1() {
	cfile := pflag.String("config",
		"config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/mysql")
	//viper.SetConfigFile("config/dev.yaml")
	//viper.KeyDelimiter("-")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
