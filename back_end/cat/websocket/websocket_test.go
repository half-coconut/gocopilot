package main

import "testing"

func TestLongConnectionRequest(t *testing.T) {
	url := "wss://api.infstones.com/ws/core/mainnet/6e97213d22994a2fae3917c0e00715d6"
	header := make(map[string][]string)
	messageType := "text"
	WSConfig := WSConfig{
		URL:                    url,
		SendMsgDurationTime:    2,
		ConnectionDurationTime: 30,
		RetryTimes:             3,
	}
	wc := NewWebsocketContent(url, header, messageType, WSConfig)
	message := "{\"jsonrpc\": \"2.0\", \"method\": \"eth_subscribe\", \"params\": [\"newHeads\"], \"id\": 1}"
	wc.LongConnectionRequest(message)
}

func TestShortConnectionRequest(t *testing.T) {
	url := "wss://api.infstones.com/ws/core/mainnet/6e97213d22994a2fae3917c0e00715d6"
	header := make(map[string][]string)
	messageType := "text"
	WSConfig := WSConfig{
		URL:                 url,
		SendMsgDurationTime: 2,
	}
	wc := NewWebsocketContent(url, header, messageType, WSConfig)
	message := "{\"jsonrpc\": \"2.0\", \"method\": \"eth_subscribe\", \"params\": [\"newHeads\"], \"id\": 1}"
	wc.ShortConnectionRequest(message)
}

func TestName(t *testing.T) {
	url := "wss://api.infstones.com/ws/core/mainnet/6e97213d22994a2fae3917c0e00715d6"
	header := make(map[string][]string)
	messageType := "text"
	WSConfig := WSConfig{
		URL:                 url,
		SendMsgDurationTime: 2,
	}
	wc := NewWebsocketContent(url, header, messageType, WSConfig)
	message := "{\"jsonrpc\": \"2.0\", \"method\": \"eth_subscribe\", \"params\": [\"newHeads\"], \"id\": 1}"
	wc.LongConnectionRequestV1(message)
}
