package domain

import (
	"time"
)

// API TODO: 待调整，抽象出来  -> 请求类型 -> 请求参数
type API struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Params  string `json:"params,omitempty"`
	Body    string `json:"body,omitempty"`
	Header  string `json:"header,omitempty"`
	Method  string `json:"method"`
	Type    string `json:"type,omitempty"` // http/websocket
	Project string `json:"project"`
	Debug   bool   `json:"debug"` // 判断是否调用接口 debug 开启 true 关闭 false

	Creator Editor `json:"creator"`
	Updater Editor `json:"updater"`
	Ctime   time.Time
	Utime   time.Time
}

type Editor struct {
	Id   int64
	Name string
}

// RecordHeader TODO: 后续把 headers 补充上
type RecordHeader struct {
	Key   []byte
	Value []byte
}
