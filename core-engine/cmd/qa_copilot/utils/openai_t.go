package utils

import (
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestEnv(t *testing.T) {
	// 加载 .env 文件
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	// 构建 .env 文件的路径
	envPath := filepath.Join(home, "Downloads", "plgo-main", "book", "aiops", "module_7", "k8scopilot", ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("err: %v", err.Error())
		log.Fatal("Error loading .env file")
	}

	a := os.Getenv("OPENAI_API_BASE")
	t.Log(a)

	// 断言环境变量不为空
	if a == "" {
		t.Error("OPENAI_API_BASE is not set")
	}

}
func TestChatGPT(t *testing.T) {

}
