package note

import (
	"context"
	"egg_yolk/internal/domain"
	"egg_yolk/internal/repository/dao/note"
	"egg_yolk/pkg/logger"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) (int64, error)
	Update(ctx context.Context, note domain.Note) error
	Sync(ctx context.Context, note domain.Note) (int64, error)
	SyncStatus(ctx context.Context, id, author_id int64, status domain.NoteStatus) error
	List(ctx context.Context, id int64, offset, limit int) ([]domain.Note, error)
}

type CacheNoteRepository struct {
	dao       note.NoteDAO
	authorDAO note.AuthorDAO
	readerDAO note.ReaderDAO
	l         logger.LoggerV1
}

func (c *CacheNoteRepository) List(ctx context.Context, id int64, offset, limit int) ([]domain.Note, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CacheNoteRepository) SyncStatus(ctx context.Context, id, author_id int64, status domain.NoteStatus) error {
	return c.dao.SyncStatus(ctx, id, author_id, status.ToUint8())
}

func NewNoteRepository(dao note.NoteDAO, authorDAO note.AuthorDAO, readerDAO note.ReaderDAO, l logger.LoggerV1) NoteRepository {
	return &CacheNoteRepository{
		dao:       dao,
		authorDAO: authorDAO,
		readerDAO: readerDAO,
		l:         l,
	}
}

func (c *CacheNoteRepository) Create(ctx context.Context, note domain.Note) (int64, error) {
	return c.dao.Insert(ctx, c.domainToEntity(note))
}

func (c *CacheNoteRepository) Update(ctx context.Context, note domain.Note) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(note))
}

func (c *CacheNoteRepository) Sync(ctx context.Context, n domain.Note) (int64, error) {
	var (
		id  = n.Id
		err error
	)
	if id > 0 {
		err = c.authorDAO.UpdateById(ctx, c.domainToEntity(n))
	} else {
		id, err = c.authorDAO.Insert(ctx, c.domainToEntity(n))
	}
	if err != nil {
		return id, err
	}
	err = c.readerDAO.Upsert(ctx, note.PublishedNote(c.domainToEntity(n)))
	return id, err
}

func (c *CacheNoteRepository) domainToEntity(n domain.Note) note.Note {
	return note.Note{
		Id:       n.Id,
		Title:    n.Title,
		Content:  n.Content,
		AuthorId: n.AuthorId,
		Role:     n.Role,
		Status:   n.Status.ToUint8(),
	}
}
