package note

import (
	"context"
	"egg_yolk/internal/domain"
)

type NoteAuthorRepository interface {
	Create(ctx context.Context, note domain.Note) (int64, error)
	Update(ctx context.Context, note domain.Note) error
}
