package main

import (
	"TestCopilot/TestEngine/cat/log"
	"fmt"
	"net/http"
	"testing"
	"time"
)

//log.InitLogger() 记得初始化log，否则空指针

func TestTimeRound(t *testing.T) {
	d1 := 1*time.Hour + 30*time.Minute + 15*time.Second
	roundedD1 := round(d1)
	fmt.Println(roundedD1) // 输出: 1h31m0s

	d2 := 1*time.Minute + 30*time.Second + 500*time.Millisecond
	roundedD2 := round(d2)
	fmt.Println(roundedD2) // 输出: 1m31s

	d3 := 1*time.Second + 500*time.Millisecond
	roundedD3 := round(d3)
	fmt.Println(roundedD3) // 输出: 2s

	d4 := 100 * time.Nanosecond
	roundedD4 := round(d4)
	fmt.Println(roundedD4) // 输出: 100ns

	d5 := 1500 * time.Millisecond
	roudedD5 := round(d5)
	fmt.Println(roudedD5)
}

func TestRadio(t *testing.T) {
	radio := 3121.0 / 3124
	fmt.Printf("%.2f%%", radio*100)
}

func token() API {
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	body := []byte(`{"email": "cc@163.com", "password": "Cc12345!"}`)
	//body := []byte(`{"email": "123@163.com", "password": "Cc12345!"}`)
	return NewAPI("login接口",
		"POST",
		"http://127.0.0.1:3002/users/login",
		"",
		"eggyolk@qq.com", body, h)
}

func user_profile() API {
	a := token()
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

func api_list() API {
	a := token()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))

	return NewAPI("API list 接口",
		"GET",
		"http://127.0.0.1:3002/api/list",
		"",
		"eggyolk@qq.com", []byte{}, h)
}

func note_publish() API {
	a := token()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))

	body := []byte(`{"id":9,"title": "Egg Yolk_9","content":"Hey, her name is Egg Yolk!","authorId":123}`)
	//body := []byte(`{"title": "Pandan_4","content":"Hey, his name is Pandan!","authorId":123}`)
	return NewAPI("API list 接口",
		"POST",
		"http://127.0.0.1:3002/note/publish",
		"",
		"eggyolk@qq.com", body, h)
}

func note_withdraw() API {
	a := token()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))

	body := []byte(`{"id":9}`)
	//body := []byte(`{"title": "Pandan_4","content":"Hey, his name is Pandan!","authorId":123}`)
	return NewAPI("API list 接口",
		"POST",
		"http://127.0.0.1:3002/note/withdraw",
		"",
		"eggyolk@qq.com", body, h)
}

func note_edit() API {
	a := token()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))

	body := []byte(`{"id":9,"title": "Egg Yolk_2","content":"Hey, his name is Pandan!","authorId":123}`)
	//body := []byte(`{"title": "Pandan_4","content":"Hey, his name is Pandan!","authorId":123}`)
	return NewAPI("API list 接口",
		"POST",
		"http://127.0.0.1:3002/note/edit",
		"",
		"eggyolk@qq.com", body, h)
}

func APIs_User() []API {
	apis := make([]API, 0)
	apis = append(apis, token(), user_profile())
	return apis
}

func APIs_API() []API {
	apis := make([]API, 0)
	apis = append(apis, token(), api_list())
	return apis
}

func APIs_Note() []API {
	apis := make([]API, 0)
	apis = append(apis, token(), note_withdraw())
	//apis = append(apis, token(), note_publish())
	//apis = append(apis, token(), note_edit())
	return apis
}

func TestApi_01_send(t *testing.T) {
	log.InitLogger()
	a := token()
	s := &subtask{
		began: time.Now(),
	}
	res := a.Send(s)
	println(res)
}

func Test_Default_ONCE(t *testing.T) {
	log.InitLogger()
	apis := APIs_Note()
	task := NewTask("默认任务，执行一次", apis, 1)
	task.DefaultRun(1, apis)
}

func Test_default_run(t *testing.T) {
	log.InitLogger()
	apis := APIs_User()
	task := NewTask("默认任务 APIs User", apis, 50)
	task.DefaultRun(2, apis)
}

func Test_http_load_profile(t *testing.T) {
	// 按照持续时间运行，以及设置 rate 运行。
	log.InitLogger()
	apis := APIs_User()
	task := NewTask("持续任务 APIs User", apis, 10)
	task.http_load(30*time.Second, 10)
}

func Test_http_load_api_list(t *testing.T) {
	// 按照持续时间运行，以及设置 rate 运行。
	log.InitLogger()
	apis := APIs_API()
	task := NewTask("持续任务 APIs API", apis, 10)
	task.http_load(30*time.Second, 10)
}
