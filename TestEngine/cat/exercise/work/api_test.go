package main

import (
	"TestCopilot/TestEngine/cat/exercise/log"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func setup_login() http.Header {
	log.InitLogger()
	// User-Agent 校验，login 和 profile 需要使用相同的value
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	body := []byte(`{"email": "cc@163.com", "password": "Cc12345!"}`)

	s := &subtask{
		began: time.Now(),
	}

	api_1 := NewAPI("login接口",
		"POST",
		"http://127.0.0.1:3002/users/login",
		"",
		"eggyolk@qq.com", body, h)
	res := api_1.Send(s)

	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))
	return h
}

func TestProfile(t *testing.T) {
	h := setup_login()
	body_2 := []byte(`{"title": "Little two","content":"Hey, his name is Xiaoer!","autherId":123}`)
	s := &subtask{
		began: time.Now(),
	}
	api_2 := NewAPI("profile接口",
		"GET",
		"http://127.0.0.1:3002/users/profile",
		"",
		"eggyolk@qq.com", body_2, h)
	api_2.Send(s)
}
