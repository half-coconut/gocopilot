package openai

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"path/filepath"
)

// OpenAI 初始化 openai client
type OpenAI struct {
	Client *openai.Client
	ctx    context.Context
}

func NewOpenAIClient() (ai *OpenAI, err error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	// 构建 .env 文件的路径
	envPath := filepath.Join(home, "Desktop", "TestCopilot", "TestEngine", "cmd", "qa_copilot", ".env")
	//envPath := filepath.Join(home, "Downloads", "TestCopilot-main", "TestEngine", "cmd", "qa_copilot", ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = os.Getenv("OPENAI_API_BASE")
	client := openai.NewClientWithConfig(config)

	ctx := context.Background()
	return &OpenAI{
		Client: client,
		ctx:    ctx,
	}, nil
}

func (o *OpenAI) SendMessage(prompt, content string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: prompt,
			}, {
				Role:    "user",
				Content: content,
			},
		},
	}
	resp, err := o.Client.CreateChatCompletion(o.ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}
	return resp.Choices[0].Message.Content, nil
}
