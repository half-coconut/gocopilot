//go:build wireinject

package main

import (
	events "TestCopilot/TestEngine/internal/events/note"
	"TestCopilot/TestEngine/internal/repository"
	"TestCopilot/TestEngine/internal/repository/cache"
	"TestCopilot/TestEngine/internal/repository/dao"
	noteDao "TestCopilot/TestEngine/internal/repository/dao/note"
	noteRepo "TestCopilot/TestEngine/internal/repository/note"
	"TestCopilot/TestEngine/internal/service"
	"TestCopilot/TestEngine/internal/service/core"
	"TestCopilot/TestEngine/internal/service/openai"
	"TestCopilot/TestEngine/internal/web"
	ijwt "TestCopilot/TestEngine/internal/web/jwt"
	"TestCopilot/TestEngine/ioc"
	"github.com/google/wire"
)

var rankingServiceSet = wire.NewSet(repository.NewCacheRankingRepository,
	cache.NewRankingRedisCache, service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, ioc.InitRLockClient,
		//ioc.InitMongoDB,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		events.NewInteractiveReadEventBatchConsumer,
		events.NewKafkaProducer,

		dao.NewUserDAO,
		dao.NewAPIDAO,
		dao.NewGORMInteractiveDAO,
		noteDao.NewNoteDAO,
		dao.NewGORMTaskDAO,

		noteDao.NewNoteAuthorDAO,
		noteDao.NewNoteReaderDAO,
		cache.NewUserCache,
		cache.NewRedisInteractiveCache,

		repository.NewUserRepository,
		noteRepo.NewNoteRepository,
		repository.NewAPIRepository,
		repository.NewCachedInteractiveRepository,
		repository.NewCacheTaskRepository,

		service.NewUserService,
		service.NewNoteService,
		service.NewAPIService,
		service.NewInteractiveService,

		core.NewTaskService,
		core.NewHttpService,
		openai.NewDeepSeekService,

		web.NewUserHandler,
		web.NewNoteHandler,
		web.NewAPIHandler,
		web.NewAIHandler,
		web.NewTaskHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddleware,
		// 组装这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
