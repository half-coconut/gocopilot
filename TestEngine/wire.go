//go:build wireinject

package main

import (
	repository2 "TestCopilot/TestEngine/interactive/repository"
	cache2 "TestCopilot/TestEngine/interactive/repository/cache"
	dao2 "TestCopilot/TestEngine/interactive/repository/dao"
	service2 "TestCopilot/TestEngine/interactive/service"
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

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var rankingServiceSet = wire.NewSet(
	repository.NewCacheRankingRepository,
	cache.NewRankingLocalCache,
	cache.NewRankingRedisCache,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, ioc.InitRLockClient,
		//ioc.InitMongoDB,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		ioc.InitIntrGRPCClient,

		interactiveSvcProvider,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		//events.NewInteractiveReadEventBatchConsumer,
		events.NewKafkaProducer,

		dao.NewUserDAO,
		dao.NewAPIDAO,
		noteDao.NewNoteDAO,
		dao.NewGORMTaskDAO,
		dao.NewGORMCronJobDAO,

		cache.NewRedisNoteCache,

		noteDao.NewNoteAuthorDAO,
		noteDao.NewNoteReaderDAO,
		cache.NewUserCache,

		repository.NewUserRepository,
		noteRepo.NewNoteRepository,
		repository.NewAPIRepository,
		repository.NewCacheTaskRepository,
		repository.NewCacheCronJobRepository,

		service.NewUserService,
		service.NewNoteService,
		service.NewAPIService,
		service.NewCronJobServiceImpl,

		core.NewTaskService,
		core.NewHttpService,
		openai.NewDeepSeekService,

		web.NewUserHandler,
		web.NewNoteHandler,
		web.NewAPIHandler,
		web.NewAIHandler,
		web.NewTaskHandler,
		web.NewCronJobHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
		ioc.InitMiddleware,
		// 组装这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
