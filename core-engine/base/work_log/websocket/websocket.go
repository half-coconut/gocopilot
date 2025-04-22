package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type WSConfig struct {
	URL                    string
	SendMsgDurationTime    time.Duration // 发送间隔时间，秒
	ConnectionDurationTime time.Duration // 连接持续时间，秒
	RetryTimes             int           // 重连次数
}

type WebsocketContent struct {
	Url         string              `json:"url"`
	Head        map[string][]string `json:"head"`
	SendMessage string              `json:"send_message"`
	MessageType string              `json:"message_type"` // "Binary"、"Text"、"Json"、"Xml"
	WSConfig    WSConfig            `json:"ws_config"`
}

func NewWebsocketContent(url string, head map[string][]string, messageType string, WSConfig WSConfig) *WebsocketContent {
	return &WebsocketContent{
		Url:         url,
		Head:        head,
		MessageType: messageType,
		WSConfig:    WSConfig,
	}
}

func (wc *WebsocketContent) ShortConnectionRequest(message string) {
	conn, _, err := websocket.DefaultDialer.Dial(wc.Url, wc.Head)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to the server!")

	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	fmt.Println("Sent:", message)

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
	}
	fmt.Printf("Received: %s\n", msg)
	fmt.Println("Connection closed.")
}

func (wc *WebsocketContent) LongConnectionRequest(message string) {
	// 创建上下文，用于控制压测持续时间
	ctx, cancel := context.WithTimeout(context.Background(), wc.WSConfig.ConnectionDurationTime*time.Second)
	defer cancel()

	// 使用 Ticker 定期发送消息
	ticker := time.NewTicker(wc.WSConfig.SendMsgDurationTime * time.Second) // 隔几秒发送一次消息
	defer ticker.Stop()

	var conn *websocket.Conn
	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if conn == nil {
				for i := 0; i < wc.WSConfig.RetryTimes; i++ {
					conn, _, err = websocket.DefaultDialer.Dial(wc.Url, wc.Head)
					if err == nil {
						fmt.Println("Connected to the server!")
						break
					}
					log.Println("Dial error:", err)
					time.Sleep(2 * time.Second) // 等待 2 秒再重试
				}

				if conn == nil {
					log.Printf("Failed to connect after retries %v times, exiting...", wc.WSConfig.RetryTimes)
					return
				}
			}
			select {
			case <-ctx.Done():
				conn.Close()
				fmt.Println("Connection closed.")
				return
			case <-ticker.C:
				// 发送和读取数据
				err = conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("Write error:", err)
					continue
				}
				fmt.Println("Sent:", message)

				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("Read error:", err)
					break
				}
				fmt.Printf("Received: %s\n", msg)
			}
		}
	}()
	wg.Wait()
}

func (wc *WebsocketContent) LongConnectionRequestV1(message string) {
	conn, _, err := websocket.DefaultDialer.Dial(wc.Url, wc.Head)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to the server!")

	// 启动一个 goroutine 用于接收消息
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			fmt.Printf("Received: %s\n", msg)
		}
	}()

	// 使用 Ticker 定期发送消息
	ticker := time.NewTicker(wc.WSConfig.SendMsgDurationTime * time.Second) // 每两秒发送一次消息
	defer ticker.Stop()

	for range ticker.C {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		fmt.Println("Sent:", message)
	}

	fmt.Println("Connection closed.")
}
