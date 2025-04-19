package note

import (
	"context"
	"time"
)

type NoteDAO interface {
	Insert(ctx context.Context, note Note) (int64, error)
	UpdateById(ctx context.Context, note Note) error
	GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Note, error)
	Sync(ctx context.Context, note Note) (int64, error)
	SyncStatus(ctx context.Context, id, authorId int64, status uint8) error
	GetById(ctx context.Context, id int64) (Note, error)
	GetPubById(ctx context.Context, id int64) (PublishedNote, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Note, error)
}
