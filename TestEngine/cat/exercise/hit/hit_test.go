package hit

import (
	"TestCopilot/TestEngine/cat/exercise/log"
	"net/http"
	"testing"
	"time"
)

func TestHit(t *testing.T) {
	log.InitLogger()
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	body := []byte(`{"email": "cc@163.com", "password": "Cc12345!"}`)

	target := NewTarget("POST", "http://127.0.0.1:3002/users/login", body, h, time.Second*30, uint64(5), uint64(10))

	target.Hitter()

}

func TestDuration(t *testing.T) {
	a := time.Duration(time.Second * 3)
	println(a)
}
