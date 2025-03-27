package execution

//
//import (
//	"TestCopilot/TestEngine/cat/log"
//	"TestCopilot/TestEngine/internal/service/core/model"
//	"net/http"
//	"sync"
//	"testing"
//	"time"
//)
//
//func api(body []byte) model.API {
//	var h = make(http.Header, 0)
//	h.Add("Content-Type", "application/json")
//	h.Add("User-Agent", "PostmanRuntime/7.39.0")
//	ht := model.NewHttpContent("POST",
//		"https://api.infstones.com/core/mainnet/6e97213d22994a2fae3917c0e00715d6",
//		"",
//		body,
//		h,
//	)
//	ws := model.WebsocketContent{}
//	api := model.NewAPI("core", "123", "ht", "egg@123.com", true, *ht, ws)
//	return api
//}
//
//func apisList() []model.API {
//	bady_1 := []byte(`{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}`)
//	bady_2 := []byte(`{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 0}`)
//	a_1 := api(bady_1)
//	a_2 := api(bady_2)
//	apis := make([]model.API, 0)
//	apis = append(apis, a_1, a_2)
//	return apis
//}
//func Test_http_load(t *testing.T) {
//	log.InitLogger()
//	apis := apisList()
//	task_conf := model.NewTaskConfig(10)
//	task := model.NewTaskService("持续任务 APIs API", apis, *task_conf)
//
//	e := NewExecutionLoadRun(task)
//	e.HttpRun(10*time.Second, 2)
//}
//
//func Test_http_debug(t *testing.T) {
//	// Debug 模式运行，这是使用的正确方式
//	log.InitLogger()
//	apis := apisList()
//	taskConf := model.NewTaskConfig(10)
//	task := model.NewTaskService("Debug任务 APIs API", apis, *taskConf)
//
//	e := NewExecutionLoadRun(task)
//
//	results := make(chan []*model.HttpResult)
//	var wg sync.WaitGroup
//	s := &model.Subtask{
//		Began: time.Now(),
//	}
//	wg.Add(1)
//	go e.HttpRunDebug(results, &wg, s)
//
//	go func() {
//		wg.Wait()
//		close(results)
//	}()
//
//	model.FinalReport(s, results)
//}
