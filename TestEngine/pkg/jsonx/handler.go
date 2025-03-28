package jsonx

import (
	"encoding/json"
	"fmt"
)

type Handler[T any] interface {
	JsonUnmarshal(jsonString string, result T)
	JsonMarshal(body T) string
}

func JsonUnmarshal[T any](jsonString string, result T) T {
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return result
}

func JsonMarshal[T any](body T) string {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	return string(jsonData)
}
