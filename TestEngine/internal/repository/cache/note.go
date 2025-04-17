package cache

import (
	"TestCopilot/TestEngine/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

type NoteCache interface {
	GetFirstPage(ctx context.Context, author int64) ([]domain.Note, error)
	SetFirstPage(ctx context.Context, author int64, notes []domain.Note) error
	DelFirstPage(ctx context.Context, author int64) error

	Set(ctx context.Context, note domain.Note) error
	Get(ctx context.Context, id int64) (domain.Note, error)

	// SetPub 正常来说，创作者和读者的 Redis 集群要分开，因为读者是一个核心中的核心
	SetPub(ctx context.Context, note domain.Note) error
	DelPub(ctx context.Context, id int64) error
	GetPub(ctx context.Context, id int64) (domain.Note, error)
}

type RedisNoteCache struct {
	client redisv9.Cmdable
}

func NewRedisNoteCache(client redisv9.Cmdable) NoteCache {
	return &RedisNoteCache{
		client: client,
	}
}

func (r *RedisNoteCache) GetPub(ctx context.Context, id int64) (domain.Note, error) {
	// 可以直接使用 Bytes 方法来获得 []byte
	data, err := r.client.Get(ctx, r.readerArtKey(id)).Bytes()
	if err != nil {
		return domain.Note{}, err
	}
	var res domain.Note
	err = json.Unmarshal(data, &res)
	return res, err
}

func (r *RedisNoteCache) SetPub(ctx context.Context, note domain.Note) error {
	data, err := json.Marshal(note)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.readerArtKey(note.Id),
		data,
		// 设置长过期时间
		time.Minute*30).Err()
}

func (r *RedisNoteCache) DelPub(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.readerArtKey(id)).Err()
}

func (r *RedisNoteCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Note, error) {
	bs, err := r.client.Get(ctx, r.firstPageKey(author)).Bytes()
	if err != nil {
		return nil, err
	}
	var notes []domain.Note
	err = json.Unmarshal(bs, &notes)
	return notes, err
}

func (r *RedisNoteCache) SetFirstPage(ctx context.Context, author int64, notes []domain.Note) error {
	// 只缓存100条
	for i := 0; i < len(notes); i++ {
		notes[i].Content = notes[i].Abstract()
	}
	data, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(author), data, time.Minute*10).Err()
}

func (r *RedisNoteCache) DelFirstPage(ctx context.Context, author int64) error {
	return r.client.Del(ctx, r.firstPageKey(author)).Err()
}

func (r *RedisNoteCache) Set(ctx context.Context, note domain.Note) error {
	data, err := json.Marshal(note)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.authorArtKey(note.Id), data, time.Minute).Err()
}

func (r *RedisNoteCache) Get(ctx context.Context, id int64) (domain.Note, error) {
	// 可以直接使用 Bytes 方法来获得 []byte
	data, err := r.client.Get(ctx, r.authorArtKey(id)).Bytes()
	if err != nil {
		return domain.Note{}, err
	}
	var res domain.Note
	err = json.Unmarshal(data, &res)
	return res, err
}

func (r *RedisNoteCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("note:first_page:%d", uid)
}

// 创作端的缓存设置
func (r *RedisNoteCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

// 读者端的缓存设置
func (r *RedisNoteCache) readerArtKey(id int64) string {
	return fmt.Sprintf("article:reader:%d", id)
}
