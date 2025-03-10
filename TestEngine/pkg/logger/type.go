package logger

// Logger 风格一
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerExample() {
	var l Logger
	phone := "188****2168"
	l.Info("用户未注册，手机号码是 %s", phone)
}

// LoggerV1 风格二: 类似 zap 的风格
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	With(args ...Field) LoggerV1
}

type Field struct {
	Key   string
	Value any
}

func LoggerV1Example() {
	var l LoggerV1
	phone := "188****2168"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}

// LoggerV2 风格三: 类似 zap 的风格
type LoggerV2 interface {
	// args 必须是偶数，并且按照 key-value，key-value 来组织
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerV2Example() {
	var l LoggerV2
	phone := "188****2168"
	l.Info("用户未注册", "phone", phone)
}

// 3种风格比较：
// Logger 兼容性最好
// LoggerV1 认同参数要有名字
// LoggerV2 有完善的代码流程，否则不建议使用
