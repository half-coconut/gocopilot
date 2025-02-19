package note

import (
	"context"
	"egg_yolk/internal/domain"
)

type NoteReaderRepository interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
}
