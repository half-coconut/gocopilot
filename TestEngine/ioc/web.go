package ioc

import (
	"TestCopilot/TestEngine/internal/web"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/internal/web/middleware"
	"TestCopilot/TestEngine/pkg/ginx/middlewares/logger"
	"TestCopilot/TestEngine/pkg/ginx/middlewares/ratelimit"
	logger2 "TestCopilot/TestEngine/pkg/logger"
	"context"
	"errors"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strings"

	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitMiddleware(redisClients redis.Cmdable, l logger2.LoggerV1, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHandler(),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("Http请求", logger2.Field{Key: "al", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePath("/hello").
			IgnorePath("/users/signup").
			IgnorePath("/users/login").Build(),
		// 限流的方案和 lua 脚本
		ratelimit.NewBuilder(redisClients, time.Second, 100).Build(),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func InitWebServer(middleware []gin.HandlerFunc, aiHandler *web.AIHandler, userHandler *web.UserHandler, apiHandler *web.APIHandler, notehandler *web.NoteHandler) *gin.Engine {
	server := gin.Default()
	server.Use(CORSMiddleware())

	server.Use(middleware...)
	parameterExamples(server)
	userHandler.RegisterRoutes(server)
	apiHandler.RegisterRoutes(server)
	aiHandler.RegisterRoutes(server)
	notehandler.RegisterRoutes(server)
	return server
}

func corsHandler() gin.HandlerFunc {
	// 使用middleware 处理跨域问题
	return cors.New(cors.Config{
		//AllowAllOrigins: true,          // 允许所有来源
		AllowedMethods: []string{"*"}, // 允许所有方法
		//AllowedHeaders: []string{"*"}, // 允许所有头
		//AllowedOrigins:  []string{"*"},
		//AllowedMethods: []string{"POST", "GET"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposedHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "test-copilot.com")
		},
		MaxAge: 12 * time.Hour,
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
