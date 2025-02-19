package service

import (
	"context"
	"egg_yolk/internal/domain"
	events "egg_yolk/internal/events/note"
	"egg_yolk/internal/repository/note"
	"egg_yolk/pkg/logger"
)

type NoteService interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
	Withdraw(ctx context.Context, note domain.Note) error
	Publish(ctx context.Context, note domain.Note) (int64, error)
	List(ctx context.Context, id int64, offset int, limit int) ([]domain.Note, error)
}

type noteService struct {
	repo     note.NoteRepository
	author   note.NoteAuthorRepository
	reader   note.NoteReaderRepository
	l        logger.LoggerV1
	producer events.Producer
}

func (n *noteService) List(ctx context.Context, id int64, offset int, limit int) ([]domain.Note, error) {
	return n.repo.List(ctx, id, offset, limit)
}

func (n *noteService) Withdraw(ctx context.Context, note domain.Note) error {
	return n.repo.SyncStatus(ctx, note.Id, note.AuthorId, domain.NoteStatusPrivate)
}

func NewNoteService(repo note.NoteRepository, l logger.LoggerV1) NoteService {
	return &noteService{
		repo: repo,
		l:    l,
	}
}

func (n *noteService) Save(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusUnpublished
	if note.Id > 0 {
		// 这里是修改
		err := n.repo.Update(ctx, note)
		if err != nil {
			n.l.Warn("修改失败", logger.Error(err))
		}
		return note.Id, err
	}
	// 这里是新增
	Id, err := n.repo.Create(ctx, note)
	if err != nil {
		n.l.Warn("新增失败", logger.Error(err))
	}
	return Id, err
}

func (n *noteService) Publish(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusPublished
	return n.repo.Sync(ctx, note)
}

func (n *noteService) PublishV1(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusPublished
	var (
		id  = note.Id
		err error
	)
	if id > 0 {
		err = n.author.Update(ctx, note)
	} else {
		id, err = n.author.Create(ctx, note)
	}
	if err != nil {
		return 0, err
	}
	note.Id = id
	for i := 0; i < 3; i++ {
		id, err = n.reader.Save(ctx, note)
		if err == nil {
			break
		}
		n.l.Error("部分失败，保存到线上库失败",
			logger.Int64("note_id", note.Id),
			logger.Error(err))
	}
	if err != nil {
		n.l.Error("部分失败，重试彻底失败",
			logger.Int64("note_id", note.Id),
			logger.Error(err))
	}
	return id, err
}
