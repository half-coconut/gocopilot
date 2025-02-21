package middleware

import (
	"TestCopilot/backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	return func(context *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}
		// 使用 JWT 来校验
		tockenHeader := context.GetHeader("Authorization")
		if tockenHeader == "" {
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "权限校验不通过"})
			return
		}
		segs := strings.Split(tockenHeader, " ")
		if len(segs) != 2 {
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "权限校验不通过"})
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "权限校验不通过"})
			return
		}

		// 对 token 进行校验
		if token == nil || !token.Valid || claims.Id == 0 {
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "权限校验不通过"})
			return
		}
		if claims.UserAgent != context.Request.UserAgent() {
			// 严重的安全问题
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "User-Agent，权限校验不通过"})
			log.Println("User-Agent，权限校验不通过")
			return
		}
		// 刷新 JWT Token
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Minute*30 {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute * 30))
			tokenStr, err = token.SignedString(web.JWTKey)
			if err != nil {
				log.Fatal("jwt 续约失败")
			}
			context.Header("x-jwt-token", tokenStr)
		}
		context.Set("claims", claims)
	}
}
