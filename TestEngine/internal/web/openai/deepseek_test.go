package openai

import (
	"testing"
)

func TestDeepseek(t *testing.T) {
	dhl := NewDeepSeekHandler()
	context, _ := dhl.DeepseekClient("你是一个资深的测试架构师", "你自我介绍一下吧")
	t.Log(context)
}
