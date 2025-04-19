package ioc

import (
	"TestCopilot/TestEngine/internal/web"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/internal/web/middleware"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/ginx/middlewares/logger"
	"TestCopilot/TestEngine/pkg/ginx/middlewares/metric"
	"TestCopilot/TestEngine/pkg/ginx/middlewares/ratelimit"
	logger2 "TestCopilot/TestEngine/pkg/logger"
	"context"
	"errors"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"

	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

func InitMiddleware(redisClients redisv9.Cmdable, l logger2.LoggerV1, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "test_copilot",
		Subsystem: "test_engine",
		Name:      "http_biz_code",
		Help:      "HTTPde 业务错误码",
	})
	return []gin.HandlerFunc{
		CORSMiddleware(),
		corsHandler(),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("Http请求", logger2.Field{Key: "al", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),

		(&metric.MiddlewareBuilder{
			Namespace:  "test_copilot",
			Subsystem:  "test_engine",
			Name:       "gin_http",
			Help:       "统计 GIN 的 HTTP 接口",
			InstanceId: "my-instance-1",
		}).Builder(),

		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePath("/hello").
			IgnorePath("/users/signup").
			IgnorePath("/users/login").
			IgnorePath("/test/metric").Build(),
		// 限流的方案和 lua 脚本，注意这里限流 200个请求
		ratelimit.NewBuilder(redisClients, time.Second, 100).Build(),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func InitWebServer(middleware []gin.HandlerFunc,
	aiHandler *web.AIHandler,
	userHandler *web.UserHandler,
	apiHandler *web.APIHandler,
	taskHandler *web.TaskHandler,
	notehandler *web.NoteHandler,
	cronJobhandler *web.CronJobHandler) *gin.Engine {
	server := gin.Default()

	server.Use(middleware...)
	parameterExamples(server)
	userHandler.RegisterRoutes(server)
	apiHandler.RegisterRoutes(server)
	taskHandler.RegisterRoutes(server)
	aiHandler.RegisterRoutes(server)
	notehandler.RegisterRoutes(server)
	cronJobhandler.RegisterRoutes(server)
	// 仅仅用于测试，不需要依赖注入
	(&web.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}

func corsHandler() gin.HandlerFunc {
	// 使用middleware 处理跨域问题
	return cors.New(cors.Config{
		AllowCredentials: true, // 允许带 cookie
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"x-jwt-token", "x-refresh-token"}, // 允许前端获取到
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		MaxAge:           12 * time.Hour,

		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			// 允许其他特定域名，比如线上域名
			if strings.Contains(origin, "test-copilot.com") {
				return true
			}
			return false
		},
	})
}

// 用于实验
func parameterExamples(server *gin.Engine) {
	// 静态路由
	server.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world")
	})
	// 参数路由
	server.GET("/users/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "这是你传过来的名字：%s", name)
	})
	// 查询参数
	server.GET("/info", func(context *gin.Context) {
		id := context.Query("id")
		context.String(http.StatusOK, "这是你传过来的 ID 是：%v", id)
	})
	// 通配符匹配
	server.GET("/views/*.html", func(context *gin.Context) {
		path := context.Param(".html")
		context.String(http.StatusOK, "匹配上的值是：%s", path)
	})
}

// 仅供实验
func initLogger() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log.Printf("这是标准日志输出")
	// 打印不出来
	zap.L().Error("之前：有问题")
	zap.ReplaceGlobals(l)
	zap.L().Error("之后：有问题")

	zap.L().Info("这是实验数据",
		zap.Error(errors.New("这个一个 error")),
		zap.String("key", "1"),
		zap.Int64("Id", 123))

}

// 仅供实验
func initLog() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 使用 go 标准输出到控制台
	log.SetOutput(os.Stdout)
}
