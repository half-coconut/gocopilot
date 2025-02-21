package note

import (
	"TestCopilot/backend/internal/domain"
	"context"
)

type NoteAuthorRepository interface {
	Create(ctx context.Context, note domain.Note) (int64, error)
	Update(ctx context.Context, note domain.Note) error
}
