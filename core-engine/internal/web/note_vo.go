package web

import "github.com/half-coconut/gocopilot/core-engine/internal/domain"

// VO view object 面向前端的
type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"` // 标记位
}

type RewardReq struct {
	Id     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

type NoteV0 struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Status   uint8  `json:"status"`

	// 计数
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	// 我个人有没有收藏，有没有点赞
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`

	Ctime string `json:"ctime"`
	Utime string `json:"utime"`
}

type NoteReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req NoteReq) toDomain(uid int64) domain.Note {
	return domain.Note{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
