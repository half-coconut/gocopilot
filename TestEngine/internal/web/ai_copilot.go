package web

import (
	"TestCopilot/TestEngine/internal/service/openai"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/pkg/ginx"
	"TestCopilot/TestEngine/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	l   logger.LoggerV1
	svc openai.DeepSeekService
}

func NewAIHandler(l logger.LoggerV1, svc openai.DeepSeekService) *AIHandler {
	return &AIHandler{
		l:   l,
		svc: svc,
	}
}

func (ai *AIHandler) RegisterRoutes(server *gin.Engine) {
	aisg := server.Group("/openai")
	aisg.POST("/ask/deepseek", ginx.WrapToken[ijwt.UserClaims](ai.AskDeepSeek))
	aisg.POST("/ask/chatgpt", ginx.WrapToken[ijwt.UserClaims](ai.AskChatGPT))
}

func (ai *AIHandler) AskChatGPT(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type AIReq struct {
		Prompt    string `json:"prompt"`
		UserInput string `json:"userInput"`
	}

	var req AIReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	// 调用 ai
	client, err := openai.NewOpenAIClient()
	if err != nil {
		return ginx.Result{Code: 0, Message: "系统错误"}, err
	}
	response, err := client.SendMessage(req.Prompt, req.UserInput)
	if err != nil {
		return ginx.Result{Code: 0, Message: "系统错误"}, err
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response,
	}, nil
}

func (ai *AIHandler) AskDeepSeek(ctx *gin.Context, uc ijwt.UserClaims) (ginx.Result, error) {
	type AIReq struct {
		Prompt    string `json:"prompt"`
		UserInput string `json:"userInput"`
	}

	var req AIReq
	err := ctx.Bind(&req)
	if err != nil {
		return ginx.Result{
			Code:    0,
			Message: "系统错误",
		}, err
	}

	// 调用 ai
	response, err := ai.svc.DeepSeekClient("你是一个资深的测试架构师", req.UserInput)
	if err != nil {
		return ginx.Result{Code: 0, Message: "系统错误"}, err
	}

	return ginx.Result{
		Code:    1,
		Message: "OK",
		Data:    response,
	}, nil
}
