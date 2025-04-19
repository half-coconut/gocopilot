package note

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
)

type NoteReaderRepository interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
}
