package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	_ "go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"net/http"
)

func main() {
	initViper()
	initPrometheus()
	//closeFunc := ioc.InitOTEL()

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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()
	//closeFunc(ctx)

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

func initViper() {
	cfile := pflag.String("config",
		"config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.SetConfigType("yaml")
	// 实时监听配置变更
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})

	fmt.Println("Config file:", *cfile) // 添加这行
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
