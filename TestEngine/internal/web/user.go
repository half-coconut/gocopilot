package web

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/service"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/logger"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserHandler struct {
	svc            service.UserService
	emailRegexp    *regexp2.Regexp
	passwordRegexp *regexp2.Regexp
	l              logger.LoggerV1
	ijwt.Handler
}

const (
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	passwordRegex = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*\W).{8,}$`
)

func NewUserHandler(svc service.UserService, l logger.LoggerV1, jwtHdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		svc:            svc,
		emailRegexp:    regexp2.MustCompile(emailRegex, regexp2.None),
		passwordRegexp: regexp2.MustCompile(passwordRegex, regexp2.None),
		l:              l,
		Handler:        jwtHdl,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组路由
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.GET("/logout", u.LogoutJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
}

func (u *UserHandler) SignUp(context *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := context.Bind(&req); err != nil {
		return
	}

	matched, err := u.emailRegexp.MatchString(req.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		u.l.Info("邮箱校验报错", logger.Error(err), logger.String("email", req.Email))
		return
	}
	if !matched {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "邮箱不正确"})
		u.l.Info("邮箱不正确", logger.Error(err), logger.String("email", req.Email))
		return
	}
	matched, err = u.passwordRegexp.MatchString(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		u.l.Info("密码校验报错", logger.Error(err), logger.String("password", req.Password))
		return
	}
	if !matched {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "密码长度不小于8位，包含数字，字母，特殊字符，字母需要大小写"})
		u.l.Info("密码校验报错", logger.Error(err), logger.String("password", req.Password))
		return
	}
	if req.Password != req.ConfirmPassword {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "两次输入密码不一致"})
		u.l.Info("两次输入密码不一致", logger.Error(err), logger.String("password", req.Password), logger.String("confirm password", req.ConfirmPassword))
		return
	}
	err = u.svc.Signup(context, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicate) {
		context.JSON(http.StatusConflict, Result{Code: 0, Message: "邮箱冲突"})
		u.l.Info("邮箱冲突", logger.Error(err), logger.String("email", req.Email))
		return
	}
	context.JSON(http.StatusOK, Result{Code: 1, Message: "注册成功！", Data: req})
	fmt.Printf("%v\n", req)

}

func (u *UserHandler) LoginSession(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.Bind(&req); err != nil {
		return
	}
	var user domain.User
	user, err := u.svc.Login(context, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "邮箱/用户或者密码不正确"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		return
	}
	// 使用 session 作为登录校验
	sess := sessions.Default(context)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		// 过期时间 30min
		MaxAge: 60 * 30,
	})
	sess.Save()
	context.JSON(http.StatusOK, Result{Code: 1, Message: "登录成功", Data: user})
	fmt.Printf("%v\n", user)
	return
}

func (u *UserHandler) LoginJWT(context *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := context.Bind(&req); err != nil {
		return
	}
	var user domain.User
	user, err := u.svc.Login(context, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "邮箱/用户或者密码不正确"})
		u.l.Info("邮箱/用户或者密码不正确", logger.Error(err), logger.String("email", req.Email))
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		u.l.Info("登录系统错误", logger.Error(err), logger.String("email", req.Email))
		return
	}
	// 使用 JWT 校验登录
	if err = u.SetLoginToken(context, user.Id); err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统异常"})
		u.l.Info("JWT登录校验，系统异常", logger.Error(err), logger.String("email", req.Email))
		return
	}

	context.JSON(http.StatusOK, Result{Code: 1, Message: "登录成功", Data: user})
	return
}

func (u *UserHandler) Logout(context *gin.Context) {
	sess := sessions.Default(context)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	context.JSON(http.StatusOK, Result{Code: 1, Message: "退出登录成功"})
	return
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code:    5,
			Message: "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Message: "退出登录OK",
	})
}

func (u *UserHandler) Edit(context *gin.Context) {
	type EditReq struct {
		Email       string `json:"email"`
		NickName    string `json:"nickName"`
		Department  string `json:"department"`
		Role        string `json:"role"`
		Description string `json:"description"`
	}
	var req EditReq
	if err := context.Bind(&req); err != nil {
		return
	}
	if req.NickName == "" {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "昵称不能为空"})
		return
	}
	if req.Department == "" {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "部门不能为空"})
		return
	}
	if req.Role == "" {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "角色不能为空"})
		return
	}
	if len(req.Description) > 1024 {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "描述过长"})
		return
	}
	err := u.svc.UpdateNonSensitiveInfo(context, domain.User{
		Email:       req.Email,
		NickName:    req.NickName,
		Department:  req.Department,
		Role:        req.Role,
		Description: req.Description,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		return
	}
	context.JSON(http.StatusOK, Result{Code: 1, Message: "更新成功", Data: req})

}

func (u *UserHandler) ProfileJWT(context *gin.Context) {
	type ProfileReq struct {
		id int64 `json:"id"`
	}
	var req ProfileReq
	err := context.Bind(&req)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	c, _ := context.Get("users")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		u.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	user, err := u.svc.Profile(context, claims.Id)

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		context.JSON(http.StatusBadRequest, Result{Code: 0, Message: "邮箱不存在"})
		u.l.Info("邮箱不存在", logger.Error(err), logger.String("email", user.Email))
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		u.l.Info("用户校验，系统错误", logger.Error(err), logger.String("email", user.Email))
		return
	}

	println(claims.Id)
	context.JSON(http.StatusOK, Result{Code: 1, Message: "获取 Profile 成功", Data: user})
}
