//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/half-coconut/gocopilot/core-engine/interactive/grpc"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository/cache"
	"github.com/half-coconut/gocopilot/core-engine/interactive/repository/dao"
	"github.com/half-coconut/gocopilot/core-engine/interactive/service"
)

var thirdProvider = wire.NewSet(InitRedis,
	InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service.NewInteractiveService(nil, nil)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
