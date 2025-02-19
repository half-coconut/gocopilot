package web

import (
	"egg_yolk/internal/domain"
	"egg_yolk/internal/service"
	"egg_yolk/pkg/ginx"
	"egg_yolk/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type NoteHandler struct {
	l   logger.LoggerV1
	svc service.NoteService
}

func NewNoteHandler(svc service.NoteService, l logger.LoggerV1) *NoteHandler {
	return &NoteHandler{
		svc: svc,
		l:   l,
	}
}

func (n *NoteHandler) RegisterRoutes(server *gin.Engine) {
	note := server.Group("/note")
	note.POST("/edit", n.Edit)
	note.POST("/publish", n.Publish)
	note.POST("/withdraw", n.Withdraw)
	note.POST("/list", ginx.WrapBodyAndToken[ListReq, UserClaims](n.List))
	note.GET("/detail:id", ginx.WrapToken[UserClaims](n.Detail))

	pub := server.Group("pub")
	pub.GET("/:id", n.PubDetail, func(ctx *gin.Context) {})
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq, UserClaims](n.Like))
	pub.POST("/reward", ginx.WrapBodyAndToken[RewardReq, UserClaims](n.Reward))
}

func (n *NoteHandler) Edit(ctx *gin.Context) {
	// 目前仅支持登录人为创建人，可新增可编辑，否则失败
	type NoteReq struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorId int64  `json:"authorId"`
		Role     string `json:"role"`
	}
	var req NoteReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	cl, _ := ctx.Get("claims")
	claims, ok := cl.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}

	note := domain.Note{
		Id:       req.Id, // 以是否传入 id，作为新增和修改的依据
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: claims.Id,
		Role:     "Author",
	}

	Id, err := n.svc.Save(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "保存成功！", Data: Id})
}

func (n *NoteHandler) Publish(ctx *gin.Context) {
	type NoteReq struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorId int64  `json:"authorId"`
		Role     string `json:"role"`
	}
	var req NoteReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	cl, _ := ctx.Get("claims")
	claims, ok := cl.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	note := domain.Note{
		Id:       req.Id, // 以是否传入 id，作为新增和修改的依据
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: claims.Id,
		Role:     "Author",
	}

	Id, err := n.svc.Publish(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("发布笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "发布成功！", Data: Id})
}

func (n *NoteHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	cl := ctx.MustGet("claims")
	claims, ok := cl.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	note := domain.Note{
		Id:       req.Id, // 以是否传入 id，作为新增和修改的依据
		AuthorId: claims.Id,
	}
	err = n.svc.Withdraw(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("撤销笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "撤销成功！"})
}

func (n *NoteHandler) List(ctx *gin.Context, req ListReq, uc UserClaims) (ginx.Result, error) {
	res, err := n.svc.List(ctx, uc.Id, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{Code: 0,
			Msg: "系统错误",
		}, nil
	}
	return ginx.Result{
		Data: slice.Map[domain.Note, NoteV0](res,
			func(idx int, src domain.Note) NoteV0 {
				return NoteV0{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					Ctime:    src.Ctime.Format(time.DateTime),
					Utime:    src.Utime.Format(time.DateTime),
				}
			}),
		Msg: "OK",
	}, nil
}

func (n *NoteHandler) Detail(ctx *gin.Context, uc UserClaims) (ginx.Result, error) {
	return ginx.Result{Msg: "OK"}, nil
}

func (n *NoteHandler) PubDetail(context *gin.Context) {
}

func (n *NoteHandler) Like(ctx *gin.Context, req LikeReq, uc UserClaims) (ginx.Result, error) {
	return ginx.Result{Msg: "OK"}, nil
}

func (n *NoteHandler) Reward(ctx *gin.Context, req RewardReq, uc UserClaims) (ginx.Result, error) {
	return ginx.Result{Msg: "OK"}, nil
}
