package ginx

import (
	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

// ginx 插件库

var L logger.LoggerV1

var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(
		opt, []string{"code"},
	)
	// 后期可以考虑 code,method,命中路由，http状态码
	prometheus.MustRegister(vector)
}

func WrapToken[C jwtv5.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cl, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := cl.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res, err := fn(ctx, c)
		if err != nil {
			L.Error("handle business logic logs",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		vector.WithLabelValues(strconv.Itoa(int(res.Code))).Inc()
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBodyAndToken[Req any, C jwtv5.Claims](fn func(ctx *gin.Context, req Req, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		err := ctx.Bind(&req)
		if err != nil {
			return
		}

		cl, ok := ctx.Get("users")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c, ok := cl.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res, err := fn(ctx, req, c)
		if err != nil {
			L.Error("handle business logic logs",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		vector.WithLabelValues(strconv.Itoa(int(res.Code))).Inc()
		ctx.JSON(http.StatusOK, res)
	}
}
