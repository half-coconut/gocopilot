package openai

import (
	httprequest "TestCopilot/TestEngine/internal/service/core/model"
	"TestCopilot/TestEngine/pkg/logger"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DeepSeekHandler struct {
}

func NewDeepSeekHandler() DeepSeekHandler {
	return DeepSeekHandler{}
}

func (d *DeepSeekHandler) DeepseekClient(prompt, userInput string) (string, error) {
	apiKey, apiEndpoint := d.getApiAndEndpoint()
	jsonBody := d.requestBody(prompt, userInput)

	var h = make(http.Header)
	h.Add("Content-Type", "application/json")
	h.Add("Authorization", "Bearer "+apiKey)

	target := httprequest.NewHttpContent("POST", apiEndpoint, "", jsonBody, h)

	s := &httprequest.Subtask{
		Began: time.Now(),
	}
	res := target.Send(s)

	var response DeepSeekResponse
	err := json.Unmarshal([]byte(res.Resp), &response)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return "", err
	}

	// 访问 content 字段
	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		log.Printf("Content: %v", content)
		return content, nil
	} else {
		log.Println("No choices found in response.")
	}
	return "", err
}

func (d *DeepSeekHandler) getApiAndEndpoint() (string, string) {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	// 构建 .env 文件的路径
	envPath := filepath.Join(home, "Desktop", "TestCopilot", "TestEngine", "pkg", "qa_copilot", ".env")
	//envPath := filepath.Join(home, "Downloads", "TestCopilot-main", "TestEngine", "pkg", "qa_copilot", ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	apiEndpoint := os.Getenv("DEEPSEEK_API_BASE")

	if apiKey == "" {
		log.Println(errors.New("DEEPSEEK_API_KEY environment variable is not set"))
	}

	return apiKey, apiEndpoint
}

func (d *DeepSeekHandler) requestBody(prompt, userInput string) []byte {
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

	//userInput := "What is the capital of France?"

	requestBody := RequestBody{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "你是一个资深的测试架构师"},
			{Role: "user", Content: userInput},
		},
		Stream: false,
	}

	// 将结构体编码为 JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal("Error marshaling JSON:", logger.Error(err))
		return nil
	}

	log.Println(string(jsonBody))
	return jsonBody
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
