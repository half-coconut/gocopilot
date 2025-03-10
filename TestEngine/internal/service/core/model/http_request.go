package model

import (
	"TestCopilot/TestEngine/pkg/logger"
	"bytes"
	"io"
	"net/http"
	"time"
)

type HttpContent struct {
	Name   string      `json:"name"`
	URL    string      `json:"url"`
	Params string      `json:"params,omitempty"`
	Body   []byte      `json:"data,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Method string      `json:"method"`
	l      logger.LoggerV1
}

func NewHttpContent(method, url, params string, body []byte, header http.Header) *HttpContent {
	return &HttpContent{
		Method: method,
		URL:    url,
		Params: params,
		Body:   body,
		Header: header,
	}
}

// HttpRequest 组合 http请求
func (h *HttpContent) HttpRequest() (*http.Request, error) {
	var body io.Reader

	if len(h.Body) != 0 {
		body = bytes.NewReader(h.Body)
	}

	req, err := http.NewRequest(h.Method, h.URL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range h.Header {
		req.Header[k] = make([]string, len(v))
		copy(req.Header[k], v)
	}

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	return req, err
}

// Send 发送 http请求
func (h *HttpContent) Send(s *Subtask) *HttpResult {
	var res HttpResult

	res.Method = h.Method
	res.URL = h.URL
	res.Req = string(h.Body)

	s.seqmu.Lock()
	res.Timestamp = s.Began.Add(time.Since(s.Began))
	res.Seq = s.seq
	s.seq++
	s.seqmu.Unlock()

	req, err := h.HttpRequest()
	if err != nil {
		h.l.Error("请求加载异常", logger.Error(err))
	}

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		h.l.Error("发送请求异常", logger.Error(err))
	}
	defer r.Body.Close()

	body := io.Reader(r.Body)
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, body); err != nil {
		h.l.Error("响应 body 复制异常", logger.Error(err))
	}
	res.Resp = string(buf.Bytes())

	if res.Code = int(uint16(r.StatusCode)); res.Code < 200 || res.Code >= 400 {
		res.Error = r.Status
	}
	res.Headers = r.Header

	defer func() {
		res.Duration = time.Since(res.Timestamp)
		if err != nil {
			res.Error = err.Error()
		}
		//h.l.Info(printResult(&res))
	}()

	return &res
}
