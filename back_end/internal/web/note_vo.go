package web

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
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
