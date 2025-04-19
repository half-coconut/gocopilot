package cmd

import (
	"github.com/half-coconut/gocopilot/core-engine/cmd/qa_scopilot/utils"
	"testing"
)

func TestEnv(t *testing.T) {
	client, _ := utils.NewOpenAIClient()
	message, err := client.SendMessage("hi", "hello")
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf(message)
}
