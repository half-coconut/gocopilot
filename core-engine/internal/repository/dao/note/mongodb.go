package note

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	//client *mongo.Client
	//db *mongo.Database

	col     *mongo.Collection // 制作库
	liveCol *mongo.Collection // 线上库
	node    *snowflake.Node
}

func (m *MongoDBDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Note, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetById(ctx context.Context, id int64) (Note, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetPubById(ctx context.Context, id int64) (PublishedNote, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Note, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) NoteDAO {
	return &MongoDBDAO{
		col:     db.Collection("notes"),
		liveCol: db.Collection("published_notes"),
		node:    node,
	}
}

func (m *MongoDBDAO) Sync(ctx context.Context, note Note) (int64, error) {
	var (
		id  = note.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, note)
	} else {
		id, err = m.Insert(ctx, note)
	}
	if err != nil {
		return 0, err
	}
	now := time.Now().UnixMilli()
	note.Utime = now
	// 操作线上库了
	filter := bson.M{"id": note.Id}
	//update := bson.M{"$set": PublishedNote(note)}
	//upsert := bson.M{"$setOnInsert": bson.M{"ctime": now}}

	updateV1 := bson.M{
		"$set":         PublishedNote(note),
		"$setOnInsert": bson.M{"ctime": now},
	}
	_, err = m.liveCol.UpdateOne(ctx, filter, updateV1,
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, note Note) error {
	filter := bson.M{"id": note.Id, "author_id": note.AuthorId}
	sets := bson.M{"$set": bson.M{
		"title":   note.Title,
		"content": note.Content,
		"status":  note.Status,
		"utime":   time.Now().UnixMilli(),
	}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败！")
	}
	return nil
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, id, authorId int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) Insert(ctx context.Context, note Note) (int64, error) {
	id := m.node.Generate().Int64()
	note.Id = id
	now := time.Now().UnixMilli()
	note.Ctime = now
	note.Utime = now
	_, err := m.col.InsertOne(ctx, note)
	return id, err
}

// ToUpdate 可以封装一些 bson 的方法，方便使用，思路：builder 模式，方便使用...
func ToUpdate(val map[string]any) bson.M {
	return val
}

func ToFilter(val map[string]any) bson.M {
	return val
}

func ToSet(val map[string]any) bson.M {
	return bson.M{"$set": val}
}

func ToUpsert(setVal, setOnInsertVal map[string]any) bson.M {
	return bson.M{
		"$set":         bson.M(setVal),
		"$setOnInsert": bson.M(setOnInsertVal),
	}
}
