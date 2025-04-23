package core

import (
	"bytes"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"io"
	"net/http"
	"time"
)

type HttpService interface {
	SetHttpInput(method, url, params string, body []byte, header http.Header)
	Send(s *Subtask) *domain.HttpResult
}
type httpService struct {
	l  logger.LoggerV1
	hc HttpContent
}

func NewHttpService(l logger.LoggerV1) HttpService {
	return &httpService{
		l: l,
	}
}

func (h *httpService) SetHttpInput(method, url, params string, body []byte, header http.Header) {
	h.hc = HttpContent{
		Method: method,
		URL:    url,
		Params: params,
		Body:   body,
		Header: header,
	}
}

func (h *httpService) Send(s *Subtask) *domain.HttpResult {
	var res domain.HttpResult

	res.Method = h.hc.Method
	res.URL = h.hc.URL
	res.Body = string(h.hc.Body)

	s.seqmu.Lock()
	res.Timestamp = s.Began.Add(time.Since(s.Began))
	res.Seq = s.seq
	s.seq++
	s.seqmu.Unlock()

	req, err := h.httpRequest()
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

	if res.Code = int64(uint16(r.StatusCode)); res.Code < 200 || res.Code >= 400 {
		res.Error = r.Status
	}
	res.Headers = r.Header

	defer func() {
		res.Duration = time.Since(res.Timestamp)
		if err != nil {
			res.Error = err.Error()
		}
	}()

	return &res
}

func (h *httpService) httpRequest() (*http.Request, error) {
	var body io.Reader

	if len(h.hc.Body) != 0 {
		body = bytes.NewReader(h.hc.Body)
	}

	req, err := http.NewRequest(h.hc.Method, h.hc.URL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range h.hc.Header {
		req.Header[k] = make([]string, len(v))
		copy(req.Header[k], v)
	}

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	return req, err
}

type HttpContent struct {
	Name   string      `json:"name"`
	URL    string      `json:"url"`
	Params string      `json:"params,omitempty"`
	Body   []byte      `json:"data,omitempty"`
	Header http.Header `json:"header,omitempty"`
	Method string      `json:"method"`
}
