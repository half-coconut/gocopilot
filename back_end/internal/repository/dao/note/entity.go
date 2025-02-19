package note

type Note struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Role     string `bson:"role,omitempty"`
	Status   uint8  `bson:"status,omitempty"`

	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}

type PublishedNote Note
