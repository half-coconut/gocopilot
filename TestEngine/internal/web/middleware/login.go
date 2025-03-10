package middleware

import (
	"TestCopilot/TestEngine/internal/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, path := range l.paths {
			if context.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(context)
		id := sess.Get("userId")
		if id == nil {
			context.JSON(http.StatusUnauthorized, web.Result{Code: 0, Message: "权限校验不通过"})
			return
		}
	}
}
