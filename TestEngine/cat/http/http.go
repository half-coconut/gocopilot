package http

import (
	"TestCopilot/TestEngine/cat/log"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type Target struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Body   []byte      `json:"body,omitempty"`
	Header http.Header `json:"header,omitempty"`
}

func NewTarget(method, url string, body []byte, header http.Header) Target {
	return Target{
		Method: method,
		URL:    url,
		Body:   body,
		Header: header,
	}
}

type Hit interface {
	Do() *Result
}

func (t *Target) request() (*http.Request, error) {
	var body io.Reader

	if len(t.Body) != 0 {
		body = bytes.NewReader(t.Body)
	}

	req, err := http.NewRequest(t.Method, t.URL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range t.Header {
		req.Header[k] = make([]string, len(v))
		copy(req.Header[k], v)
	}

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}
	return req, err
}

func (t *Target) Do() *Result {
	var res Result

	res.Method = t.Method
	res.URL = t.URL

	start := time.Now()

	req, err := t.request()
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
		res.Duration = time.Since(start).String()
		if err != nil {
			res.Error = err.Error()
		}
		log.L.Info(printResult(&res))
	}()

	return &res
}

func printResult(res *Result) string {
	return fmt.Sprintf("Code: %d\nError: %s\nBody: %s\nMethod: %s\nURL:%s\nHeaders: %v\nDuration: %v\n", res.Code, res.Error, res.Body, res.Method, res.URL, res.Headers, res.Duration)
}

type Result struct {
	Code     int         `json:"code,string"`
	Error    string      `json:"error"`
	Body     string      `json:"body"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Headers  http.Header `json:"headers"`
	Duration string      `json:"duration"`
}
