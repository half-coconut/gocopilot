package service

import (
	"TestCopilot/TestEngine/internal/domain"
	events "TestCopilot/TestEngine/internal/events/note"
	"TestCopilot/TestEngine/internal/repository/note"
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"time"
)

//go:generate mockgen -source=note.go -package=mocks -destination=mocks/note.mock.go NoteService
type NoteService interface {
	Save(ctx context.Context, note domain.Note) (int64, error)
	Withdraw(ctx context.Context, note domain.Note) error
	Publish(ctx context.Context, note domain.Note) (int64, error)
	List(ctx context.Context, id int64, offset int, limit int) ([]domain.Note, error)
	ListPub(ctx context.Context, offset int, limit int) ([]domain.Note, error)
	GetById(ctx context.Context, id int64) (domain.Note, error)
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Note, error)
}

type noteService struct {
	repo     note.NoteRepository
	author   note.NoteAuthorRepository
	reader   note.NoteReaderRepository
	l        logger.LoggerV1
	producer events.Producer

	ch chan readInfo
}

func (svc *noteService) ListPub(ctx context.Context, offset int, limit int) ([]domain.Note, error) {
	//TODO implement me
	panic("implement me")
}

type readInfo struct {
	uid int64
	nid int64
}

func (svc *noteService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Note, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	note, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		// 每次打开一篇文章，就发一条消息
		go func() {
			// 生产者也可以通过改批量来提高性能
			er := svc.producer.ProducerReadEvent(
				ctx,
				events.ReadEvent{
					// 即便你的消费者要用 art 的里面的数据，
					// 让它去查询，你不要在 event 里面带
					// 除非是快照语义的可以加 content
					Uid: uid,
					Nid: id,
				})
			if er != nil {
				svc.l.Error("发送读者阅读事件失败")
			}
		}()

		go func() {
			// 改批量的做法
			svc.ch <- readInfo{
				uid: uid,
				nid: id,
			}
		}()

	}
	return note, err
}

func (svc *noteService) GetById(ctx context.Context, id int64) (domain.Note, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *noteService) List(ctx context.Context, id int64, offset int, limit int) ([]domain.Note, error) {
	return svc.repo.List(ctx, id, offset, limit)
}

func (svc *noteService) Withdraw(ctx context.Context, note domain.Note) error {
	return svc.repo.SyncStatus(ctx, note.Id, note.Author.Id, domain.NoteStatusPrivate)
}

func NewNoteService(repo note.NoteRepository, l logger.LoggerV1, producer events.Producer) NoteService {
	return &noteService{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}

func NewNoteServiceV2(repo note.NoteRepository, l logger.LoggerV1, producer events.Producer) NoteService {
	ch := make(chan readInfo, 10)
	go func() {
		for {
			uids := make([]int64, 0, 10)
			nids := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			for i := 0; i < 10; i++ {
				select {
				case info, ok := <-ch:
					if !ok {
						cancel()
						return
					}
					uids = append(uids, info.uid)
					nids = append(nids, info.uid)
				case <-ctx.Done():
					break

				}
			}
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			producer.ProducerReadEventV1(ctx, events.ReadEventV1{
				Uids: uids,
				Nids: nids,
			})
			cancel()
		}
	}()
	return &noteService{
		repo:     repo,
		l:        l,
		producer: producer,
		ch:       ch,
	}
}

func (svc *noteService) Save(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusUnpublished
	if note.Id > 0 {
		// 这里是修改
		err := svc.repo.Update(ctx, note)
		if err != nil {
			svc.l.Warn("修改失败", logger.Error(err))
		}
		return note.Id, err
	}
	// 这里是新增
	Id, err := svc.repo.Create(ctx, note)
	if err != nil {
		svc.l.Warn("新增失败", logger.Error(err))
	}
	return Id, err
}

func (svc *noteService) Publish(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusPublished
	return svc.repo.Sync(ctx, note)
}

func (svc *noteService) PublishV1(ctx context.Context, note domain.Note) (int64, error) {
	note.Status = domain.NoteStatusPublished
	var (
		id  = note.Id
		err error
	)
	if id > 0 {
		err = svc.author.Update(ctx, note)
	} else {
		id, err = svc.author.Create(ctx, note)
	}
	if err != nil {
		return 0, err
	}
	note.Id = id
	for i := 0; i < 3; i++ {
		id, err = svc.reader.Save(ctx, note)
		if err == nil {
			break
		}
		svc.l.Error("部分失败，保存到线上库失败",
			logger.Int64("note_id", note.Id),
			logger.Error(err))
	}
	if err != nil {
		svc.l.Error("部分失败，重试彻底失败",
			logger.Int64("note_id", note.Id),
			logger.Error(err))
	}
	return id, err
}
