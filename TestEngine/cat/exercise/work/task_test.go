package main

import (
	"TestCopilot/TestEngine/cat/exercise/log"
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"time"
)

// log.InitLogger() 记得初始化log，否则空指针

func api_1() API {
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	body := []byte(`{"email": "cc@163.com", "password": "Cc12345!"}`)
	return NewAPI("login接口",
		"POST",
		"http://127.0.0.1:3002/users/login",
		"",
		"eggyolk@qq.com", body, h)
}

func api_2() API {
	a := api_1()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))

	body := []byte(`{"title": "Little two","content":"Hey, his name is Xiaoer!","autherId":123}`)
	return NewAPI("profile接口",
		"GET",
		"http://127.0.0.1:3002/users/profile",
		"",
		"eggyolk@qq.com", body, h)
}

func APIs() []API {
	apis := make([]API, 0)
	apis = append(apis, api_1(), api_2())
	return apis
}

func TestApi_01_send(t *testing.T) {
	log.InitLogger()
	a := api_1()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	println(res)
}

func TestDisplayTask(t *testing.T) {
	log.InitLogger()
	apis := APIs()
	task := NewTask("profile接口", apis, 100)
	println(task)
	fmt.Printf("%v", task)
	log.L.Info(displayTask(task))
}

func TestWorkRun(t *testing.T) {
	log.InitLogger()
	apis := APIs()
	fmt.Println("当前 goroutine 数量:", runtime.NumGoroutine())
	WorkRun(2, "Work任务", apis)
}

func TestDefaultRun(t *testing.T) {
	log.InitLogger()
	apis := APIs()
	task := NewTask("默认任务", apis, 50)
	task.DefaultRun(2, apis)
}

func TestConstantDefaultRun(t *testing.T) {
	log.InitLogger()
	apis := APIs()
	task := NewTask("持续任务", apis, 10)
	task.ConstantDefaultRun(2, apis, 2*time.Second)
}

func Test_http_load(t *testing.T) {
	// 按照持续时间运行，以及设置 rate 运行。
	log.InitLogger()
	apis := APIs()
	task := NewTask("持续任务", apis, 10)
	task.http_load(10*time.Second, 2)
}
