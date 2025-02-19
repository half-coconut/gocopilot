package ioc

import (
	"context"
	"egg_yolk/internal/web"
	"egg_yolk/internal/web/middleware"
	"egg_yolk/pkg/ginx/middlewares/logger"
	"egg_yolk/pkg/ginx/middlewares/ratelimit"
	logger2 "egg_yolk/pkg/logger"
	"errors"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitMiddleware(redisClients redis.Cmdable, l logger2.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHandler(),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("Http请求", logger2.Field{Key: "al", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePath("/hello").
			IgnorePath("/users/signup").
			IgnorePath("/users/login").Build(),
		// 限流的方案和 lua 脚本
		ratelimit.NewBuilder(redisClients, time.Second, 100).Build(),
	}
}

func InitWebServer(middleware []gin.HandlerFunc, userHandler *web.UserHandler, apiHandler *web.APIHandler, notehandler *web.NoteHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middleware...)
	parameterExamples(server)
	userHandler.RegisterRoutes(server)
	apiHandler.RegisterRoutes(server)
	notehandler.RegisterRoutes(server)
	return server
}

func corsHandler() gin.HandlerFunc {
	// 使用middleware 处理跨域问题
	return cors.New(cors.Config{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		// 为了使用 jwt
		ExposedHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "egg-yolk.com")
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
