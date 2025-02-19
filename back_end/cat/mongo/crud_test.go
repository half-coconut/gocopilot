package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("egg_yolk")
	col := mdb.Collection("notes")

	//defer func() {
	//	_, err = col.DeleteOne(ctx, bson.D{})
	//}()
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()

	// 新增数据
	now := time.Now().UnixMilli()
	res, err := col.InsertOne(ctx, Note{
		Id:       id,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 4,
		Role:     "author",
		Ctime:    now,
		Utime:    now,
	})
	assert.NoError(t, err)
	// _id mongo 的 id
	fmt.Printf("id: %s \n", res.InsertedID)

	// 查询数据
	filter := bson.D{bson.E{"id", id}}
	var note Note
	err = col.FindOne(ctx, filter).Decode(&note)
	assert.NoError(t, err)
	fmt.Printf("%v \n", note)

	// 实践中，多用这个写法，同时，注意定义好 omitempty
	note = Note{}
	err = col.FindOne(ctx, Note{Id: id}).Decode(&note)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("没有数据")
	} else {
		assert.NoError(t, err)
	}
	fmt.Printf("%v \n", note)

	// 查询数据 or
	or := bson.A{bson.M{"id": 125}, bson.M{"id": 123}}
	orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	assert.NoError(t, err)
	var notes []Note
	err = orRes.All(ctx, &notes)
	assert.NoError(t, err)
	fmt.Printf("+++++or 查询结果+++++：%v \n", notes)

	// 查询数据 and
	and := bson.A{bson.M{"title": "我的标题"}, bson.M{"id": 123}}
	andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	assert.NoError(t, err)
	notes = []Note{}
	err = andRes.All(ctx, &notes)
	assert.NoError(t, err)
	fmt.Printf("+++++and 查询结果+++++：%v \n", notes)

	// 查询数据 in
	in := bson.D{bson.E{"id", bson.M{"$in": []any{123, 125}}}}
	//in := bson.D{bson.E{"id", bson.D{bson.E{"$in", []any{123, 125}}}}}
	inRes, err := col.Find(ctx, in)
	assert.NoError(t, err)
	notes = []Note{}
	err = inRes.All(ctx, &notes)
	assert.NoError(t, err)
	fmt.Printf("+++++in 查询结果+++++：%v \n", notes)

	// options 查询，查询特定字段
	opsRes, err := col.Find(ctx, in, options.Find().SetProjection(
		bson.M{"id": 1, "title": 1}))
	assert.NoError(t, err)
	notes = []Note{}
	err = opsRes.All(ctx, &notes)
	assert.NoError(t, err)
	fmt.Printf("+++++options 查询结果+++++：%v \n", notes)

	// 更新数据
	//filter_update := bson.D{bson.E{Key: "id", Value: 125}}
	filter_update := bson.M{"id": id}
	sets := bson.D{bson.E{Key: "$set", Value: bson.M{"title": "新的标题-new"}}}
	updateRes, err := col.UpdateMany(ctx, filter_update, sets)
	assert.NoError(t, err)
	fmt.Println("affect: ", updateRes.ModifiedCount)

	// 更新数据
	filter_update = bson.M{"id": id}
	updateRes, err = col.UpdateMany(ctx, filter_update, bson.D{bson.E{Key: "$set", Value: Note{AuthorId: 2, Role: "editor"}}})
	assert.NoError(t, err)
	fmt.Println("affected: ", updateRes.ModifiedCount)

	// 创建索引
	//col.Indexes().CreateOne(ctx, mongo.IndexModel{
	//	Keys:    bson.M{"id": 1},
	//	Options: options.Index().SetUnique(true),
	//})

	// 创建多个索引
	//idxRes, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
	//	{Keys: bson.M{"id": 1},
	//		Options: options.Index().SetUnique(true)},
	//	{Keys: bson.M{"author_id": 1}},
	//})
	//assert.NoError(t, err)
	//fmt.Printf("+++++index+++++：%v \n", idxRes)

	// 删除数据
	filter_del := bson.D{bson.E{Key: "id", Value: 124}}
	deleteRes, err := col.DeleteOne(ctx, filter_del)
	assert.NoError(t, err)
	fmt.Println("deleted: ", deleteRes.DeletedCount)

}

//type Note struct {
//	Id       int64  `bson:"id"`
//	Title    string `bson:"title"`
//	Content  string `bson:"content"`
//	AuthorId int64  `bson:"author_id"`
//	Role     string `bson:"role"`
//	Status   uint8  `bson:"status"`
//
//	Ctime int64 `bson:"ctime"`
//	Utime int64 `bson:"utime"`
//}

type Note struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Role     string `bson:"role,omitempty"`
	Status   uint8  `bson:"status,omitempty"`

	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}
