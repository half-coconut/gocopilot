package jwt

import (
	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractToken(ctx *gin.Context) string
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwtv5.RegisteredClaims
}

type UserClaims struct {
	jwtv5.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Id        int64  `json:"id"`
	Ssid      string `json:"ssid"`
	UserAgent string `json:"user_agent"`
	VIP       bool   `json:"vip"`
}
