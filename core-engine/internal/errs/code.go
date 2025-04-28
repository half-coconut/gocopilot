package errs

const (
	// CommonInvalidInput 任何模块都可以使用
	CommonInvalidInput        = 400001
	CommonInternalServerError = 500001
)

// 用户模块
const (
	// UserInvalidInput 用户模块输入错误，含糊的错误
	UserInvalidInput        = 401001
	UserInternalServerError = 501001
	// UserInvalidOrPassword 用户不存在，或者密码错误
	UserInvalidOrPassword = 401002
)

// Note 模块
const (
	NoteInvalidInput        = 402001
	NoteInternalServerError = 502001
)

// 任务模块
const (
	TaskInvalidInput        = 403001
	TaskInternalServerError = 503001
	TaskInvalidOrNotExist   = 403002
)

// CronJob模块
const (
	CronJobInvalidInput        = 404001
	CronJobInternalServerError = 504001
	CronJobInvalidOrNotExist   = 404002
)

var (
	UserInvalidInputV1 = Code{
		Number: 401001,
		Msg:    "用户输入错误",
	}
)

type Code struct {
	Number int
	Msg    string
}
