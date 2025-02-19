package domain

import "time"

type Note struct {
	Id       int64 // 用于鉴别是新增还是修改
	Title    string
	Content  string
	AuthorId int64
	Role     string
	Status   NoteStatus

	Ctime time.Time
	Utime time.Time
}

type Role struct {
	Author string // 首次创建，默认为 Auther，自己编辑时
	Editor string // 其他人编辑时
	Viewer string // get 请求时，仅查看
}

// PublishedArticle 衍生类型，偷个懒
type PublishedNote Note

type NoteStatus uint8

func (n Note) Abstract() string {
	// 将 a.Content 转换为 rune 切片。这是为了正确处理中文字符，
	// 因为在 Go 中，一个中文字符通常占用多个字节。
	// 如果直接使用 string 的切片，会导致中文字符被拆分，可能出现乱码。
	cs := []rune(n.Content)
	if len(cs) < 100 {
		return n.Content
	}
	return string(cs[:100])
}

const (
	// NoteStatusUnknown 为了避免零值的问题
	NoteStatusUnknown NoteStatus = iota
	NoteStatusUnpublished
	NoteStatusPublished
	NoteStatusPrivate
)

func (a NoteStatus) ToUint8() uint8 {
	return uint8(a)
}

func (a NoteStatus) Valid() bool {
	return a.ToUint8() > 0
}

func (a NoteStatus) NonPublished() bool {
	return a != NoteStatusPublished
}

func (a NoteStatus) ToString() string {
	switch a {
	case NoteStatusUnpublished:
		return "unpublished"
	case NoteStatusPublished:
		return "published"
	case NoteStatusPrivate:
		return "private"
	default:
		return "unknown"
	}
}

// NoteStatusV1 如果状态很复杂，有很多行为，定义很多方法，或者有很多额外的字段
type NoteStatusV1 struct {
	Val  uint8
	Name string
}

var (
	NoteStatusV1Unknown = NoteStatusV1{
		Val: 0, Name: "unknown",
	}
)
