package note

import (
	"TestCopilot/TestEngine/internal/domain"
	"context"
)

type NoteReaderRepository interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
}
