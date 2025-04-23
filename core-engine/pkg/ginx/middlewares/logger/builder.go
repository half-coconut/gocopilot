package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method:  ctx.Request.Method,
			URL:     url,
			Headers: ctx.Request.Header,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			// Body 是个数据流，读完就没了
			body, _ := ctx.GetRawData()
			// ctx.GetRawData() 内部就是调用 io.ReadAll()
			// body, _ := io.ReadAll(ctx.Body.Body)
			// 读完后还要方回去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

			// 1024 可以作为参数传入
			if len(body) > 1024 {
				body = body[:1024]
			}
			// 这是一个很消耗 CPU 和内存的操作，因为会引起复制
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				ResponseWriter: ctx.Writer,
				al:             al,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()
		// 执行到业务逻辑
		ctx.Next()

	}
}

// responseWriter 这是 gin 的扩展库
type responseWriter struct {
	// 组合的装饰器模式，这是特定的几个方法
	gin.ResponseWriter
	al *AccessLog
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	// 要写回去
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	//if len(data) > 1024 {
	//	data = data[:1024]
	//}
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	//if len(data) > 1024 {
	//	data = data[:1024]
	//}
	w.al.RespBody = string(data)
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// Http 请求的方法
	Method   string
	URL      string
	Status   int
	Headers  http.Header
	ReqBody  string
	RespBody string
	Duration string
}
