package ioc

import (
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

var mongoDB *mongo.Database

func InitMongoDB() *mongo.Database {
	if mongoDB == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context,
				startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
		}
		opts := options.Client().
			ApplyURI("mongodb://root:example@localhost:27017/").
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mongoDB = client.Database("testengine")
	}
	return mongoDB
}

func InitDB(l logger.LoggerV1) *gorm.DB {
	// 使用 gorm 打印日志
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/testengine"), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有查询时间超过这个阈值，才会使用
			// 50ms, 100ms
			// SQL 查询要求命中索引，最好走一次磁盘IO，不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true, // 线上环境设置为 true，比较好
			ParameterizedQueries:      false,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
