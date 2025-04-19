//go:build wireinject

package main

import (
	events "github.com/half-coconut/gocopilot/core-engine/interactive/events"
	"github.com/half-coconut/gocopilot/core-engine/interactive/grpc"
	"github.com/half-coconut/gocopilot/core-engine/interactive/ioc"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/interactive/service"

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
