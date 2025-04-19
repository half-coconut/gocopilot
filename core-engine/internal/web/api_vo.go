package web

type APIListResponse struct {
	Interfaces []API0 `json:"interfaces"` // API 列表
	Total      int    `json:"total"`      // API 总数
}

// 前端得到的API数据
type API0 struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Params      string `json:"params"`
	Body        string `json:"body"`
	Header      string `json:"header"`
	Method      string `json:"method"`
	Type        string `json:"type"` // http/websocket
	Project     string `json:"project"`
	DebugResult string `json:"debug_result"`

	Creator string `json:"creator"`
	Updater string `json:"updater"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

type APIReq struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Params  string `json:"params"`
	Type    string `json:"type"`
	Body    string `json:"body"`
	Header  string `json:"header"`
	Method  string `json:"method"`
	Project string `json:"project"`
	Debug   bool   `json:"debug"`
}

var (
	Body   map[string]interface{}
	Header map[string]string
)
