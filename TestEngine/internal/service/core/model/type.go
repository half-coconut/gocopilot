package model

// API 接口结构体，API -> taskService
type API struct {
	Name     string           `json:"name"`
	TeamId   string           `json:"team_id"`
	LoadType string           `json:"load_type"` // 接口类型：http,websocket
	Http     HttpContent      `json:"http_content"`
	WS       WebsocketContent `json:"ws_content"`
	Debug    bool             `json:"debug"`
	Creator  string           `json:"creator"`
	Updater  string           `json:"updater"`
}

func NewAPI(name, teamId, types, email string, debug bool, http HttpContent, ws WebsocketContent) API {
	return API{
		Name:     name,
		TeamId:   teamId,
		LoadType: types,
		Http:     http,
		WS:       ws,
		Debug:    debug,
		Creator:  email,
		Updater:  email,
	}
}

// Assert TODO:简单断言
type Assert struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
