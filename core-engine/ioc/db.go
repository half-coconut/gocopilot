package ioc

import (
	"context"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/config"
	dao2 "github.com/half-coconut/gocopilot/core-engine/interactive/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	promesdk "github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

var mongoDB *mongo.Database

func InitMongoDB() *mongo.Database {
	if mongoDB == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				fmt.Println(startedEvent.Command)
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

			},
		}
		opts := options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mongoDB = client.Database("coreengine")
	}
	return mongoDB
}

func InitDB(l logger.LoggerV1) *gorm.DB {
	// 使用 gorm 打印日志
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有查询时间超过这个阈值，才会使用
			// 50ms, 100ms
			// SQL 查询要求命中索引，最好走一次磁盘IO，不到 100ms
			SlowThreshold:             time.Millisecond * 100,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true, // 线上环境设置为 true，比较好
			ParameterizedQueries:      false,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "coreengine",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	// 监控查询的执行时间
	// 如果是 JOIN 查询，table 就是 JOIN 在一起
	// 或者 table 就是主表，A JOIN B，记录的就是 A
	pcb := newCallbacks()
	//pcb.registerAll(db)
	db.Use(pcb)

	tracing.NewPlugin(tracing.WithDBSystem("coreengine"), tracing.WithQueryFormatter(func(query string) string {
		l.Debug("", logger.String("query", query))
		return query
	}),
	//tracing.WithoutMetrics(), // 不要记录 metrics 各种指标
	//tracing.WithoutQueryVariables(), // 不要记录查询参数，线上
	)

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	err = dao2.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

type Callbacks struct {
	vector *promesdk.SummaryVec
}

func (pcb *Callbacks) Name() string {
	return "prometheus-query"
}

func (pcb *Callbacks) Initialize(db *gorm.DB) error {
	pcb.registerAll(db)
	return nil
}

func newCallbacks() *Callbacks {
	vector := promesdk.NewSummaryVec(promesdk.SummaryOpts{
		// 设施各种 namespace
		Namespace: "go_copilot",
		Subsystem: "core_engine",
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的 执行时间",
		ConstLabels: map[string]string{
			"db": "coreengine",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})

	pcb := &Callbacks{
		vector: vector,
	}
	promesdk.MustRegister(vector)
	return pcb
}

func (pcb *Callbacks) registerAll(db *gorm.DB) {
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", pcb.after("create"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", pcb.after("update"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", pcb.after("delete"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", pcb.after("raw"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", pcb.after("row"))
	if err != nil {
		panic(err)
	}
}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}
func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}
		// 准备上报 prometheus
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).Observe(float64(time.Since(startTime).Milliseconds()))
	}
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
