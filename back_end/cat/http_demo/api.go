package main

import (
	"bytes"
	"egg_yolk/cat/log"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

// API 接口结构体，API -> Task 之间，应该还是有个 Case
// API 这里默认是 http 请求，后期可以补充，例如 websocket RPC等
type API struct {
	Name    string      `json:"name"`
	URL     string      `json:"url"`
	Params  string      `json:"params,omitempty"`
	Body    []byte      `json:"data,omitempty"`
	Header  http.Header `json:"header,omitempty"`
	Method  string      `json:"method"`
	Creator string      `json:"creator"`
	Updater string      `json:"updater"`
}

func NewAPI(name, method, url, params, email string, body []byte, header http.Header) API {
	return API{
		Name:    name,
		Method:  method,
		URL:     url,
		Params:  params,
		Body:    body,
		Header:  header,
		Creator: email,
		Updater: email,
	}
}

// HttpRequest 组合 http请求
func (a *API) HttpRequest() (*http.Request, error) {
	var body io.Reader

	if len(a.Body) != 0 {
		body = bytes.NewReader(a.Body)
	}

	req, err := http.NewRequest(a.Method, a.URL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range a.Header {
		req.Header[k] = make([]string, len(v))
		copy(req.Header[k], v)
	}

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	return req, err
}

// Send 发送 http请求
func (a *API) Send(s *subtask) *Result {
	var res Result

	res.Method = a.Method
	res.URL = a.URL

	s.seqmu.Lock()
	res.Timestamp = s.began.Add(time.Since(s.began))
	res.Seq = s.seq
	s.seq++
	s.seqmu.Unlock()

	req, err := a.HttpRequest()
	if err != nil {
		log.L.Error("请求加载异常", zap.Error(err))
	}

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		log.L.Error("发送请求异常", zap.Error(err))
	}
	defer r.Body.Close()

	body := io.Reader(r.Body)
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, body); err != nil {
		log.L.Error("响应 body 复制异常", zap.Error(err))
	}
	res.Body = string(buf.Bytes())

	if res.Code = int(uint16(r.StatusCode)); res.Code < 200 || res.Code >= 400 {
		res.Error = r.Status
	}
	res.Headers = r.Header

	defer func() {
		res.Duration = time.Since(res.Timestamp)
		if err != nil {
			res.Error = err.Error()
		}
		//log.L.Info(printResult(&res))
	}()

	return &res
}

func printResult(res *Result) string {
	return fmt.Sprintf(`
+++++ API Result: +++++
[Code: %d]
[Method: %s]
[URL:%s]
[Duration: %v]
[Headers: %v]
[Body: %s]
[Error: %s]`, res.Code, res.Method, res.URL, res.Duration, res.Headers, res.Body, res.Error)
}
