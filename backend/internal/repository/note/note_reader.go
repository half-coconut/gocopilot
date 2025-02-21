package note

import (
	"TestCopilot/backend/internal/domain"
	"context"
)

type NoteReaderRepository interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
}
