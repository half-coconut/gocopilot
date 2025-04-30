package note

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/cache"
	dao "github.com/half-coconut/gocopilot/core-engine/internal/repository/dao/note"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"time"
)

type NoteRepository interface {
	Create(ctx context.Context, note domain.Note) (int64, error)
	Update(ctx context.Context, note domain.Note) error
	Sync(ctx context.Context, note domain.Note) (int64, error)
	SyncStatus(ctx context.Context, id, author_id int64, status domain.NoteStatus) error
	List(ctx context.Context, uid int64, offset, limit int) ([]domain.Note, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Note, error)
	GetByID(ctx context.Context, id int64) (domain.Note, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Note, error)
}

type CacheNoteRepository struct {
	dao dao.NoteDAO
	// 这里直接在 repository 层面上操作 user 的操作
	// 如果在微服务架构下，可以直接使用用户服务来实现 UserRepository，提供用户数据
	userRepo  repository.UserRepository
	cache     cache.NoteCache
	authorDAO dao.AuthorDAO
	readerDAO dao.ReaderDAO
	l         logger.LoggerV1
}

func NewNoteRepository(dao dao.NoteDAO, userRepo repository.UserRepository, authorDAO dao.AuthorDAO, readerDAO dao.ReaderDAO, cache cache.NoteCache, l logger.LoggerV1) NoteRepository {
	return &CacheNoteRepository{
		dao:       dao,
		userRepo:  userRepo,
		authorDAO: authorDAO,
		readerDAO: readerDAO,
		cache:     cache,
		l:         l,
	}
}

func (c *CacheNoteRepository) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Note, error) {
	res, err := c.dao.ListPub(ctx, start, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(res, func(idx int, src dao.Note) domain.Note {
		return c.entityToDomain(src)
	}), nil
}

func (c *CacheNoteRepository) GetPublishedById(ctx context.Context, id int64) (domain.Note, error) {
	// 读取线上库数据
	note, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Note{}, nil
	}
	// 适合单体应用
	user, err := c.userRepo.FindById(ctx, note.AuthorId)
	if err != nil {
		return domain.Note{}, nil
	}
	res := domain.Note{
		Id:      note.Id, // 用于鉴别是新增还是修改
		Title:   note.Title,
		Content: note.Content,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.FullName},
		Status: domain.NoteStatus(note.Status),
		Ctime:  time.UnixMilli(note.Ctime),
		Utime:  time.UnixMilli(note.Utime),
	}
	return res, nil
}

func (c *CacheNoteRepository) GetByID(ctx context.Context, id int64) (domain.Note, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Note{}, err
	}
	return c.entityToDomain(data), nil
}

func (c *CacheNoteRepository) List(ctx context.Context, uid int64, offset, limit int) ([]domain.Note, error) {
	// 在这里集成缓存方案
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			// 注意：是否为同步或者异步的调用，最好有调用者来决定，通常下层就只提供一个同步的方法。
			go func() {
				c.preCache(ctx, data)
			}()
			return data, err
		}
	}
	res, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}

	data, err := slice.Map[dao.Note, domain.Note](res, func(idx int, src dao.Note) domain.Note {
		return c.entityToDomain(src)
	}), nil
	// 回写缓存，如果并发不高，直接 set
	// 如果有很高并发，就 del
	// 可以同步，也可以异步

	go func() {
		err := c.cache.SetFirstPage(ctx, uid, data)
		c.l.Error("回写缓存失败", logger.Error(err))
		c.preCache(ctx, data)
	}()
	return data, nil
}

func (c *CacheNoteRepository) SyncStatus(ctx context.Context, id, author_id int64, status domain.NoteStatus) error {
	return c.dao.SyncStatus(ctx, id, author_id, status.ToUint8())
}

func (c *CacheNoteRepository) Create(ctx context.Context, note domain.Note) (int64, error) {
	id, err := c.dao.Insert(ctx, c.domainToEntity(note))
	if err == nil {
		err = c.cache.DelFirstPage(ctx, note.Author.Id)
		if err != nil {
			c.l.Error(fmt.Sprintf("删除缓存失败：%v", err))
		}
		err = c.cache.SetPub(ctx, note)
		if err != nil {
			c.l.Error(fmt.Sprintf("set缓存失败：%v", err))
		}
	}
	return id, err
}

func (c *CacheNoteRepository) Update(ctx context.Context, note domain.Note) error {
	defer func() {
		// 如果有数据变更，清空缓存
		err := c.cache.DelFirstPage(ctx, note.Author.Id)
		if err != nil {
			return
		}
	}()
	return c.dao.UpdateById(ctx, c.domainToEntity(note))
}

func (c *CacheNoteRepository) Sync(ctx context.Context, n domain.Note) (int64, error) {
	defer func() {
		// 如果有数据变更，清空缓存
		err := c.cache.DelFirstPage(ctx, n.Author.Id)
		if err != nil {
			return
		}
	}()

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
	err = c.readerDAO.Upsert(ctx, dao.PublishedNote(c.domainToEntity(n)))
	return id, err
}

func (c *CacheNoteRepository) domainToEntity(n domain.Note) dao.Note {
	return dao.Note{
		Id:       n.Id,
		Title:    n.Title,
		Content:  n.Content,
		AuthorId: n.Author.Id,
		Status:   n.Status.ToUint8(),
	}
}

func (u *CacheNoteRepository) entityToDomain(n dao.Note) domain.Note {
	return domain.Note{
		Id:      n.Id, // 用于鉴别是新增还是修改
		Title:   n.Title,
		Content: n.Content,
		Author:  domain.Author{Id: n.AuthorId},
		Status:  domain.NoteStatus(n.Status),
		Ctime:   time.UnixMilli(n.Ctime),
		Utime:   time.UnixMilli(n.Utime),
	}
}

func (c *CacheNoteRepository) preCache(ctx context.Context, data []domain.Note) {
	// 缓存预加载方案
	// 简单处理缓存大对象的问题 1024*1024 超过 1MB 不缓存
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := c.cache.Set(ctx, data[0])
		if err != nil {
			c.l.Error("提前预加载缓存失败", logger.Error(err))
		}
	}
}
