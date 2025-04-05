//go:build wireinject

package startup

import (
	"TestCopilot/TestEngine/interactive/repository"
	"TestCopilot/TestEngine/interactive/repository/cache"
	"TestCopilot/TestEngine/interactive/repository/dao"
	"TestCopilot/TestEngine/interactive/service"
	"github.com/google/wire"
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

//func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
//	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
//	return new(grpc.InteractiveServiceServer)
//}
