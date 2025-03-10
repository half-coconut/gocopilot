package ginx

type Result struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}
