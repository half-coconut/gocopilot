package model

import (
	"TestCopilot/TestEngine/cat/exercise/log"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	log.InitLogger()
	s := &Subtask{
		Began: time.Now(),
	}
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.39.0")
	body := []byte(`{"email": "test@123.com", "password": "Cc12345!"}`)

	target := NewHttpContent("POST", "http://127.0.0.1:3002/users/login", "", body, h)
	res := target.Send(s)
	t.Log(res)
}

func TestCoreRequest(t *testing.T) {
	log.InitLogger()
	url := "https://api.infstones.com/neo/mainnet/cdfeb4ccba2b4b7faab8178d77c09788"
	//url := "https://api.infstones.com/core/mainnet/6e97213d22994a2fae3917c0e00715d6"
	//url := "https://api.infstones.com/ethereum/mainnet/f6ee2ff15aa64aedac9dde1bfc7ca45f"

	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("User-Agent", "PostmanRuntime/7.43.0")
	h.Add("Cookie", "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE2ODYxNDcsImlhdCI6MTcxMTY3ODk0NywiZGF0YSI6eyJjb21wb25lbnQiOiJ3ZWJzZXJ2aWNlIiwiZW1haWwiOiI1ODQ5NDc1NTlAcXEuY29tIiwidWlkIjoiODMyOTAwNDIiLCJyb2xlIjoiYWRtaW4iLCJnaWQiOiIzOTUzNTI0NSIsIm1mYV9lbmFibGVkIjpmYWxzZX19.duVCGmKOaYcJVoDdvklDFMYEfNx1ofBA8brc5YYz2IQ")

	body := []byte(`{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}`)

	ht := NewHttpContent("POST",
		url,
		"",
		body,
		h,
	)

	s := &Subtask{
		Began: time.Now(),
	}
	res := ht.Send(s)
	t.Log(res)
}

func TestHttpRequest(t *testing.T) {
	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	//h.Add("User-Agent", "PostmanRuntime/7.39.0")
	h.Add("Authorization", "Bearer "+"sk-c6f66dd1346a40d3b5f61bbc43aa1695")
	body := []byte(`{
        "model": "deepseek-chat",
        "messages": [
          {"role": "system", "content": "You are a helpful assistant."},
          {"role": "user", "content": "Hello!"}
        ],
        "stream": false
      }`)

	target := NewHttpContent("POST", "https://api.deepseek.com/chat/completions", "", body, h)

	s := &Subtask{
		Began: time.Now(),
	}
	res := target.Send(s)

	t.Log(res)
}

func TestDeepseekV2(t *testing.T) {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// 定义请求体结构体
	type RequestBody struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}

	userContent := "What is the capital of France?"

	// 构造请求体
	requestBody := RequestBody{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: userContent},
		},
		Stream: false,
	}

	// 将结构体编码为 JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Error("Error marshaling JSON:", err)
		return
	}

	// 打印 JSON 数据 (可选)
	t.Log(string(jsonBody))

	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("Authorization", "Bearer "+"sk-c6f66dd1346a40d3b5f61bbc43aa1695")

	target := NewHttpContent("POST", "https://api.deepseek.com/chat/completions", "", jsonBody, h)

	s := &Subtask{
		Began: time.Now(),
	}
	res := target.Send(s)

	t.Log(res.Resp)

	var response DeepSeekResponse
	err = json.Unmarshal([]byte(res.Resp), &response)
	if err != nil {
		t.Error("Error decoding JSON:", err)
		return
	}

	// 访问 content 字段
	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		t.Log("Content:", content)
	} else {
		t.Error("No choices found in response.")
	}
}

// 定义 Usage 结构体
type Usage struct {
	PromptTokens          int            `json:"prompt_tokens"`
	CompletionTokens      int            `json:"completion_tokens"`
	TotalTokens           int            `json:"total_tokens"`
	PromptTokensDetails   map[string]int `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int            `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int            `json:"prompt_cache_miss_tokens"`
}

// 定义 Message 结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义 Choice 结构体
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	LogProbs     interface{} `json:"logprobs"` // 可以是 null， 所以用 interface{}
	FinishReason string      `json:"finish_reason"`
}

// 定义 顶层结构体
type DeepSeekResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}
