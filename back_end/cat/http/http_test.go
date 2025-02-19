package http

import (
	"egg_yolk/cat/log"
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

	target := NewTarget("POST", "http://127.0.0.1:3002/users/login", body, h)
	res := target.Do()

	h.Add("Authorization", fmt.Sprintf("Bearer %s", res.Headers["X-Jwt-Token"][0]))
	return h
}

func TestLogin(t *testing.T) {
	h := setup_login()
	target := NewTarget("GET", "http://127.0.0.1:3002/users/profile", []byte{}, h)
	target.Do()
}

func TestNote_add_new_note_withAuthId(t *testing.T) {
	h := setup_login()
	body := []byte(`{"title": "Little two","content":"Hey, his name is Xiaoer!","autherId":123}`)
	target := NewTarget("POST", "http://127.0.0.1:3002/note/edit", body, h)
	target.Do()
}

func TestNote_add_new_note_withoutAuthId(t *testing.T) {
	h := setup_login()
	bady := []byte(`{"title": "Little two","content":"Hey, his name is Xiaoer!"}`)
	target := NewTarget("POST", "http://127.0.0.1:3002/note/edit", bady, h)
	target.Do()
}

func TestNote_updat_old_note_withAuthId(t *testing.T) {
	h := setup_login()
	bady := []byte(`{"id":3,"title": "Pandan2","content":"Hey, his name is Pandan!","autherId":123}`)
	target := NewTarget("POST", "http://127.0.0.1:3002/note/edit", bady, h)
	target.Do()
}

func Test_http_load(t *testing.T) {
	h := setup_login()
	target := NewTarget("GET", "http://127.0.0.1:3002/users/profile", []byte{}, h)
	target.http_load(3*time.Second, 5)
}
