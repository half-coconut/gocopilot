package hit

import (
	"TestCopilot/TestEngine/base/exercise/log"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

type Target struct {
	Method     string        `json:"method"`
	URL        string        `json:"url"`
	Body       []byte        `json:"body,omitempty"`
	Header     http.Header   `json:"header,omitempty"`
	Duration   time.Duration `json:"duration"`
	workers    uint64
	maxWorkers uint64
}

func NewTarget(method, url string, body []byte, header http.Header, Duration time.Duration, workers, maxWorkers uint64) Target {
	return Target{
		Method:     method,
		URL:        url,
		Body:       body,
		Header:     header,
		Duration:   Duration,
		workers:    workers,
		maxWorkers: maxWorkers,
	}
}

type Hit interface {
	Hitter()
	Do() *Result
}

func (t *Target) Hitter() {
	var wg sync.WaitGroup
	log.L.Debug(fmt.Sprintf("%s", t.URL))
	workers := t.workers

	if workers > t.maxWorkers {
		workers = t.maxWorkers
	}
	// 记录结果的 channel
	results := make(chan *Result)

	wg.Add(int(t.maxWorkers))

	for i := uint64(0); i < workers; i++ {
		log.L.Debug(fmt.Sprintf("i: %d", i))
		go t.attack(&wg, results)

	}

	go func() {
		wg.Wait()
		close(results)
	}()

	t.Display(results)

}

func (t *Target) attack(wg *sync.WaitGroup, results chan<- *Result) {
	wg.Done()
	results <- t.Do()
	log.L.Debug("t.Do 写入 result")
	//select {
	//case results <- t.Do():
	//case <-time.After(time.Second * 5): // 设置 5 秒超时
	//	log.L.InterfacesDebug("attack timeout")
	//}
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
	}()
	log.L.Debug("这是来自 Do 的日志：" + printResult(&res))

	return &res
}

func printResult(res *Result) string {
	return fmt.Sprintf("Code: %d\nError: %s\nBody: %s\nMethod: %s\nURL:%s\nHeaders: %v\nDuration: %v\n", res.Code, res.Error, res.Body, res.Method, res.URL, res.Headers, res.Duration)
}

func (t *Target) Display(results chan *Result) {
	for res := range results {
		log.L.Debug(fmt.Sprintf("Code: %d\nError: %s\nBody: %s\nMethod: %s\nURL:%s\nHeaders: %v\nDuration: %v\n", res.Code, res.Error, res.Body, res.Method, res.URL, res.Headers, res.Duration))
	}
}
