//go:build wireinject

package main

import (
	"egg_yolk/internal/repository"
	"egg_yolk/internal/repository/cache"
	"egg_yolk/internal/repository/dao"
	"egg_yolk/internal/repository/dao/note"
	note2 "egg_yolk/internal/repository/note"
	"egg_yolk/internal/service"
	"egg_yolk/internal/web"
	"egg_yolk/ioc"
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
