//go:build wireinject

package main

import (
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/internal/repository/cache"
	"TestCopilot/TestEngine/internal/repository/dao"
	"TestCopilot/TestEngine/internal/repository/dao/note"
	note2 "TestCopilot/TestEngine/internal/repository/note"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/web"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		//ioc.InitMongoDB,
		ioc.InitLogger,

		dao.NewUserDAO,
		note.NewNoteDAO,
		dao.NewAPIDAO,
		dao.NewGORMInteractiveDAO,
		note.NewNoteAuthorDAO,
		note.NewNoteReaderDAO,
		cache.NewUserCache,
		cache.NewRedisInteractiveCache,
		repository.NewUserRepository,
		note2.NewNoteRepository,
		repository.NewAPIRepository,
		repository.NewCachedInteractiveRepository,
		service.NewUserService,
		service.NewNoteService,
		service.NewAPIService,
		service.NewInteractiveService,
		web.NewUserHandler,
		web.NewNoteHandler,
		web.NewAPIHandler,
		web.NewAIHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddleware,
	)
	return new(gin.Engine)
}
