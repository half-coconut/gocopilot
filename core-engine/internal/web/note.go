package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	intrv1 "github.com/half-coconut/gocopilot/core-engine/api/proto/gen/intr/v1"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/service"
	ijwt "github.com/half-coconut/gocopilot/core-engine/internal/web/jwt"
	"github.com/half-coconut/gocopilot/core-engine/pkg/ginx"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*NoteHandler)(nil)

type NoteHandler struct {
	l       logger.LoggerV1
	svc     service.NoteService
	intrSvc intrv1.InteractiveServiceClient
	biz     string
}

func NewNoteHandler(svc service.NoteService, l logger.LoggerV1, intrSvc intrv1.InteractiveServiceClient) *NoteHandler {
	return &NoteHandler{
		svc:     svc,
		l:       l,
		intrSvc: intrSvc,
		biz:     "note",
	}
}

func (n *NoteHandler) RegisterRoutes(server *gin.Engine) {
	note := server.Group("/note")
	note.POST("/edit", n.Edit)
	note.POST("/publish", n.Publish)
	note.POST("/withdraw", n.Withdraw)
	note.POST("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](n.List))
	note.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](n.Detail))

	pub := note.Group("/pub")
	pub.GET("/:id", n.PubDetail, func(ctx *gin.Context) {})
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq, ijwt.UserClaims](n.Like))
	pub.POST("/reward", ginx.WrapBodyAndToken[RewardReq, ijwt.UserClaims](n.Reward))
}

func (n *NoteHandler) Edit(ctx *gin.Context) {
	// 目前仅支持登录人为创建人，可新增可编辑，否则失败
	var req NoteReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	c, _ := ctx.Get("users")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}

	note := req.toDomain(claims.Id)

	Id, err := n.svc.Save(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("保存笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "保存成功！", Data: Id})
}

func (n *NoteHandler) Publish(ctx *gin.Context) {
	var req NoteReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	c, _ := ctx.Get("users")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}

	note := req.toDomain(claims.Id)

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
	c, _ := ctx.Get("users")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("未发现用户 token 信息：%v", claims.Id), logger.Error(err))
		return
	}
	note := domain.Note{
		Id: req.Id, // 以是否传入 id，作为新增和修改的依据
	}
	err = n.svc.Withdraw(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result{Code: 0, Message: "系统错误"})
		n.l.Info(fmt.Sprintf("撤销笔记失败，用户 Id：%v", claims.Id), logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Message: "撤销成功！"})
}

func (n *NoteHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (ginx.Result, error) {
	res, err := n.svc.List(ctx, uc.Id, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{Code: 0,
			Message: "系统错误",
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
		Message: "OK",
	}, nil
}

func (n *NoteHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return ginx.Result{
			Code:    4,
			Message: "参数错误",
		}, err
	}
	note, err := n.svc.GetById(ctx, id)
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("获得文章信息失败", logger.Error(err))
		return ginx.Result{
			Code:    5,
			Message: "系统错误",
		}, err
	}
	// 这是不借助数据库查询来判定的方法
	if note.Author.Id != uc.Id {
		//ctx.JSON(http.StatusOK)
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		//a.l.Error("非法访问文章，创作者 ID 不匹配",
		//	logger.Int64("uid", usr.Id))
		return ginx.Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Message: "输入有误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", uc.Id)
	}
	return ginx.Result{
		Data: NoteV0{
			Id:    note.Id,
			Title: note.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  note.Status.ToUint8(),
			Content: note.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: note.Ctime.Format(time.DateTime),
			Utime: note.Utime.Format(time.DateTime),
		},
	}, nil
}

func (n *NoteHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code:    4,
			Message: "参数错误",
		})
		n.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}

	// 读文章本体
	uc := ctx.MustGet("users").(ijwt.UserClaims)

	// 对 waitGroup 的二次开发，在业务开发中常见
	var eg errgroup.Group
	var note domain.Note
	eg.Go(func() error {
		note, err = n.svc.GetPublishedById(ctx, id, uc.Id)
		if err != nil {
			n.l.Error("获取 publicedBy uid 失败", logger.Error(err))
		}
		return err
	})

	// 获得这篇文章的全部计数
	var getResp *intrv1.GetResponse
	eg.Go(func() error {
		getResp, err = n.intrSvc.Get(ctx, &intrv1.GetRequest{
			Biz:   n.biz,
			BizId: id,
			Uid:   uc.Id,
		})
		return err
	})

	// 前两个error 的返回，注意返回了 Error 是不容错的写法
	err = eg.Wait()
	if err != nil {
		// 查询出错
		ctx.JSON(http.StatusOK, Result{
			Code:    5,
			Message: "系统错误",
		})
		return
	}

	// 增加阅读计数
	go func() {
		// 最好不要在gorutine 里面复用外面的 error
		_, er := n.intrSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
			Biz:   n.biz,
			BizId: id,
		})
		if er != nil {
			n.l.Error("增加阅读计数失败", logger.Int64("noteId: ", note.Id), logger.Error(er))
		}
	}()
	intr := getResp.Intr

	ctx.JSON(http.StatusOK, Result{
		Data: NoteV0{
			Id:         note.Id,
			Title:      note.Title,
			Status:     note.Status.ToUint8(),
			Content:    note.Content,
			Author:     note.Author.Name,
			Ctime:      note.Ctime.Format(time.DateTime),
			Utime:      note.Utime.Format(time.DateTime),
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			LikeCnt:    intr.LikeCnt,
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
		},
	})
}

func (n *NoteHandler) Like(ctx *gin.Context, req LikeReq, uc ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		_, err = n.intrSvc.Like(ctx, &intrv1.LikeRequest{
			Biz:   n.biz,
			BizId: req.Id,
			Uid:   uc.Id,
		})
	} else {
		_, err = n.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz:   n.biz,
			BizId: req.Id,
			Uid:   uc.Id,
		})
	}

	if err != nil {
		return ginx.Result{
			Code:    5,
			Message: "系统错误",
		}, err
	}
	return ginx.Result{Message: "OK"}, nil
}

func (n *NoteHandler) Reward(ctx *gin.Context, req RewardReq, uc ijwt.UserClaims) (ginx.Result, error) {
	return ginx.Result{Message: "OK"}, nil
}
