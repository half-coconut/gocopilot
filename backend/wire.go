//go:build wireinject

package main

import (
	"TestCopilot/backend/internal/repository"
	"TestCopilot/backend/internal/repository/cache"
	"TestCopilot/backend/internal/repository/dao"
	"TestCopilot/backend/internal/repository/dao/note"
	note2 "TestCopilot/backend/internal/repository/note"
	"TestCopilot/backend/internal/service"
	"TestCopilot/backend/internal/web"
	"TestCopilot/backend/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewUserDAO,
		note.NewNoteDAO,
		dao.NewAPIDAO,
		note.NewNoteAuthorDAO,
		note.NewNoteReaderDAO,
		cache.NewUserCache,
		repository.NewUserRepository,
		note2.NewNoteRepository,
		repository.NewAPIRepository,
		service.NewUserService,
		service.NewNoteService,
		service.NewAPIService,
		web.NewUserHandler,
		web.NewNoteHandler,
		web.NewAPIHandler,

		ioc.InitWebServer,
		ioc.InitMiddleware,
	)
	return new(gin.Engine)
}
