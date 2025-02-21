package cmd

import (
	"k8scopilot/utils"
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
