package domain

import (
	"net/http"
	"time"
)

// API TODO: 待调整，抽象出来  -> 请求类型 -> 请求参数
type API struct {
	Id      int64       `json:"id"`
	Name    string      `json:"name"`
	URL     string      `json:"url"`
	Params  string      `json:"params,omitempty"`
	Body    string      `json:"data,omitempty"`
	Header  http.Header `json:"header,omitempty"`
	Method  string      `json:"method"`
	Type    string      `json:"type"` // http/websocket
	Project string      `json:"project"`

	Creator int64 `json:"creator"`
	Updater int64 `json:"updater"`
	Ctime   time.Time
	Utime   time.Time
}
