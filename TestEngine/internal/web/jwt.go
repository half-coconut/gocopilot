package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Id        int64
	UserAgent string
}

// JWTKey 因为 JWT Key 不太可能变，所以可以直接写成常量
// 也可以考虑做成依赖注入
var JWTKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")

func (u *UserHandler) SetJWTToken(context *gin.Context, Id int64) (string, error) {
	claims := UserClaims{
		Id: Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: context.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}
	context.Header("x-jwt-token", tokenStr)
	fmt.Printf("JWT Token: %v\n", tokenStr)
	return tokenStr, nil
}
