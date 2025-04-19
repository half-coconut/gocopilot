package note

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
)

type NoteAuthorRepository interface {
	Create(ctx context.Context, note domain.Note) (int64, error)
	Update(ctx context.Context, note domain.Note) error
}
