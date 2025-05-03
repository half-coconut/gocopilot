package dao

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ReportDAO interface {
	InsertDebugLog(ctx context.Context, log DebugLog) (int64, error)
	InsertSummary(ctx context.Context, s Summary) (int64, error)
	GetDebugLogsByTaskId(ctx context.Context, tid int64) ([]DebugLog, error)
}

type MongoDBReportDAO struct {
	//client    *mongo.Client
	db        *mongo.Database
	l         logger.LoggerV1
	node      *snowflake.Node
	debugLogs *mongo.Collection
	summary   *mongo.Collection
}

func NewMongoDBReportDAO(db *mongo.Database, l logger.LoggerV1) ReportDAO {
	node, _ := snowflake.NewNode(1) // 线上环境从环境变量中获取
	return &MongoDBReportDAO{
		db:        db,
		debugLogs: db.Collection("debug_logs"),
		summary:   db.Collection("summary"),
		node:      node,
		l:         l,
	}
}

func (m *MongoDBReportDAO) InsertDebugLog(ctx context.Context, log DebugLog) (int64, error) {
	id := m.node.Generate().Int64()
	log.Id = id
	now := time.Now().UnixMilli()
	log.Ctime = now
	log.Utime = now
	_, err := m.debugLogs.InsertOne(ctx, log)
	return id, err
}

func (m *MongoDBReportDAO) InsertSummary(ctx context.Context, s Summary) (int64, error) {
	id := m.node.Generate().Int64()
	s.Id = id
	now := time.Now().UnixMilli()
	s.Ctime = now
	s.Utime = now
	_, err := m.summary.InsertOne(ctx, s)
	return id, err
}

func (m *MongoDBReportDAO) GetDebugLogsByTaskId(ctx context.Context, tid int64) ([]DebugLog, error) {
	//TODO implement me
	panic("implement me")
}

type DebugLog struct {
	Id       int64       `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	TaskId   int64       `gorm:"index" bson:"task_id,omitempty"`
	AId      int64       `gorm:"index" bson:"a_id,omitempty"`
	AName    string      `bson:"a_name,omitempty"` // 接口名称
	Request  RequestInfo `bson:"request,omitempty"`
	Response interface{} `bson:"response,omitempty"`
	ClientIP string      `bson:"client_ip,omitempty"`
	Error    string      ` bson:"error,omitempty"`

	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}

type RequestInfo struct {
	URL     string            `bson:"url,omitempty"`
	Method  string            `bson:"method,omitempty"`
	Headers map[string]string `bson:"headers,omitempty"`
	Body    interface{}       `bson:"body,omitempty"`
}

type Summary struct {
	Id            int64      `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"` // 使用雪花算法生成的字符串ID
	TaskId        int64      `gorm:"index" bson:"task_id,omitempty"`
	AIds          int64      `gorm:"index" bson:"a_ids,omitempty"`
	TName         string     `bson:"t_name,omitempty"` // 任务名称 (使用string, 空字符串表示null)
	Debug         bool       `bson:"debug"`            // 开启或者关闭 Debug 模式
	Total         int        `bson:"total,omitempty"`
	Rate          float64    `bson:"rate,omitempty"`
	Throughput    float64    `bson:"throughput,omitempty"`
	TotalDuration int64      `bson:"total_duration,omitempty"`
	Min           int64      `bson:"min,omitempty"`
	Mean          int64      `bson:"mean,omitempty"`
	Max           int64      `bson:"max,omitempty"`
	P50           int64      `bson:"p50,omitempty"`
	P90           int64      `bson:"p90,omitempty"`
	P95           int64      `bson:"p95,omitempty"`
	P99           int64      `bson:"p99,omitempty"`
	Ratio         float64    `bson:"ratio,omitempty"`
	StatusCodes   string     `bson:"status_codes,omitempty"` // 存储状态码统计信息
	TestStatus    TestStatus `bson:"test_status,omitempty"`

	Status int   `bson:"status,omitempty"`
	Ctime  int64 `bson:"ctime,omitempty"`
	Utime  int64 `bson:"utime,omitempty"`
}

// TestStatus 统计个数
type TestStatus struct {
	Passed  int64
	Failed  int64
	Skipped int64
	Errors  int64
}

const (
	Unknown = iota
	Passed
	Failed
	Skipped
	Errors
)
