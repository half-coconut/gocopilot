package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

//go:embed hello.py
var helloPy string

func main() {
	// 创建临时文件
	tmpFile, err := ioutil.TempFile("", "hello*.py")
	if err != nil {
		log.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name()) // 确保程序退出时删除临时文件

	// 将嵌入的 Python 代码写入临时文件
	if _, err := tmpFile.WriteString(helloPy); err != nil {
		log.Fatalf("Error writing to temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("Error closing temp file: %s", err)
	}

	// 执行 Python 脚本
	cmd := exec.Command("python3", tmpFile.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing Python script: %s\nOutput: %s", err, string(out))
	}

	fmt.Printf("Python script output:\n%s\n", string(out))
}
