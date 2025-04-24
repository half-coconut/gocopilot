//go:build wireinject

package main

import (
	"github.com/google/wire"
	repository2 "github.com/half-coconut/gocopilot/core-engine/interactive/repository"
	cache2 "github.com/half-coconut/gocopilot/core-engine/interactive/repository/cache"
	dao2 "github.com/half-coconut/gocopilot/core-engine/interactive/repository/dao"
	service2 "github.com/half-coconut/gocopilot/core-engine/interactive/service"
	events "github.com/half-coconut/gocopilot/core-engine/internal/events/note"
	reportEvt "github.com/half-coconut/gocopilot/core-engine/internal/events/report"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao"
	noteDao "github.com/half-coconut/gocopilot/core-engine/internal/repository/dao/note"
	noteRepo "github.com/half-coconut/gocopilot/core-engine/internal/repository/note"
	"github.com/half-coconut/gocopilot/core-engine/internal/service"
	"github.com/half-coconut/gocopilot/core-engine/internal/service/core"
	"github.com/half-coconut/gocopilot/core-engine/internal/service/openai"
	"github.com/half-coconut/gocopilot/core-engine/internal/web"
	ijwt "github.com/half-coconut/gocopilot/core-engine/internal/web/jwt"
	"github.com/half-coconut/gocopilot/core-engine/ioc"
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
		ioc.InitMongoDB,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		//ioc.InitIntrGRPCClient,
		// 放一起，启用了 etcd 作为配置中心
		ioc.InitEtcd,
		ioc.InitIntrGRPCClientV1,

		// 这是流量控制用的
		//interactiveSvcProvider,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		//events.NewInteractiveReadEventBatchConsumer,
		events.NewKafkaNoteProducer,
		reportEvt.NewKafkaDebugLogProducer,
		reportEvt.NewDebugLogEventConsumer,

		dao.NewUserDAO,
		dao.NewAPIDAO,
		noteDao.NewNoteDAO,
		dao.NewGORMTaskDAO,
		dao.NewGORMCronJobDAO,
		dao.NewMongoDBReportDAO,

		cache.NewRedisNoteCache,

		noteDao.NewNoteAuthorDAO,
		noteDao.NewNoteReaderDAO,
		cache.NewUserCache,

		repository.NewUserRepository,
		noteRepo.NewNoteRepository,
		repository.NewAPIRepository,
		repository.NewCacheTaskRepository,
		repository.NewCacheCronJobRepository,
		repository.NewUncachedReportRepository,

		service.NewUserService,
		service.NewNoteService,
		service.NewAPIService,
		service.NewCronJobService,

		core.NewReportService,
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
