//go:build wireinject

package main

import (
	events "TestCopilot/TestEngine/interactive/events"
	"TestCopilot/TestEngine/interactive/grpc"
	"TestCopilot/TestEngine/interactive/ioc"
	"TestCopilot/TestEngine/interactive/repository"
	"TestCopilot/TestEngine/interactive/repository/cache"
	"TestCopilot/TestEngine/interactive/repository/dao"
	"TestCopilot/TestEngine/interactive/service"

	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitAPP() *App {
	wire.Build(
		interactiveSvcProvider,
		thirdPartySet,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
