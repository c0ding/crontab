package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

// 任务的执行时间点
type TimePoint struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// 一条日志
type LogRecord struct {
	JobName   string    `json:"job_name"`   // 任务名
	Command   string    `json:"command"`    // shell命令
	Err       string    `json:"err"`        // 脚本错误
	Content   string    `json:"content"`    // 脚本输出
	TimePoint TimePoint `json:"time_point"` // 执行时间点
}

// jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"` // JobName赋值为job10
}

// startTime小于某时间
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp} }
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

var (
	err        error
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	result     *mongo.InsertOneResult
	resultMany *mongo.InsertManyResult
	cursor     mongo.Cursor
	record     *LogRecord
	docId      objectid.ObjectID
	logArr     []interface{}

	delCond   *DeleteCond
	delResult *mongo.DeleteResult

	//cond *FindByJobName
)

func connect() {
	if client, err = mongo.Connect(context.TODO(), "mongodb://localhost:27017",
		clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
	}
	database = client.Database("day1")
	collection = database.Collection("aaa")
}

func main() {

	connect()
	//add()
	findAll()

}

func findAll() {

	// 4, 按照jobName字段过滤, 想找出jobName=job10, 找出5条
	//cond = &FindByJobName{JobName: "jobd哈哈111"} // {"jobName": "job10"}

	// 5, 查询（过滤 +翻页参数）

	if cursor, err = collection.Find(context.TODO(), nil); err != nil {
		fmt.Println(err)
		return
	}

	// 延迟释放游标
	defer cursor.Close(context.TODO())

	// 6, 遍历结果集
	for cursor.Next(context.TODO()) {
		// 定义一个日志对象
		record = &LogRecord{}

		// 反序列化bson到对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		// 把日志行打印出来
		fmt.Println(*record)
	}
}

func update() {

}
func add() {
	record = &LogRecord{
		"job10",
		"echo hello",
		"error",
		"hello",
		TimePoint{StartTime: time.Now().Unix(),
			EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		goto ERR
	}

	docId = result.InsertedID.(objectid.ObjectID)
	fmt.Println("自增ID:", docId.Hex(), result.InsertedID) //ObjectID  相等

	logArr = []interface{}{record, record, record}
	if resultMany, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		goto ERR
	}

	for _, insertID := range resultMany.InsertedIDs {
		// 拿着interface{}， 反射成objectID
		docId = insertID.(objectid.ObjectID)
		fmt.Println("222:=== ", docId.Hex())
	}

ERR:
	fmt.Println(err)
}

func delete() {
	// 4, 要删除开始时间早于当前时间的所有日志($lt是less than)
	//  delete({"timePoint.startTime": {"$lt": 当前时间}})
	delCond = &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	// 执行删除
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}
}
