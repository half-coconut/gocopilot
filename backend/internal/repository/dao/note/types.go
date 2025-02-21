package note

import "context"

type NoteDAO interface {
	Insert(ctx context.Context, note Note) (int64, error)
	UpdateById(ctx context.Context, note Note) error
	Sync(ctx context.Context, note Note) (int64, error)
	SyncStatus(ctx context.Context, id, authorId int64, status uint8) error
}
