package data

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

const (
	logRepoQueryTimeout  = 10 * time.Second
	logRepoAggTimeout    = 15 * time.Second
	logRepoInsertTimeout = 5 * time.Second
)

// LogRepository 日志仓库接口
type LogRepository interface {
	QueryLogs(ctx context.Context, filter bson.M, page, pageSize int) ([]bson.M, int64, error)
	InsertLog(ctx context.Context, log interface{}) error
	DeleteLogs(ctx context.Context, filter bson.M) (int64, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error)
}

type MongoLogRepo struct {
	col *mongo.Collection
}

func NewLogRepository(col *mongo.Collection) LogRepository {
	return &MongoLogRepo{col: col}
}

func (r *MongoLogRepo) QueryLogs(ctx context.Context, filter bson.M, page, pageSize int) ([]bson.M, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, logRepoQueryTimeout)
	defer cancel()
	findOpts := options.Find().SetSort(bson.M{"request_time": -1}).SetSkip(int64((page-1)*pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		log.Printf("[MongoLogRepo.QueryLogs] 查询失败: %v\n", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	for cursor.Next(ctx) {
		var m bson.M
		if err := cursor.Decode(&m); err == nil {
			results = append(results, m)
		}
	}
	total, _ := r.col.CountDocuments(ctx, filter)
	return results, total, nil
}

func (r *MongoLogRepo) InsertLog(ctx context.Context, logData interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, logRepoInsertTimeout)
	defer cancel()
	_, err := r.col.InsertOne(ctx, logData)
	if err != nil {
		log.Printf("[MongoLogRepo.InsertLog] 插入失败: %v\n", err)
	}
	return err
}

func (r *MongoLogRepo) DeleteLogs(ctx context.Context, filter bson.M) (int64, error) {
	res, err := r.col.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *MongoLogRepo) Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(ctx, logRepoAggTimeout)
	defer cancel()
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("[MongoLogRepo.Aggregate] 聚合失败: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	for cursor.Next(ctx) {
		var m bson.M
		if err := cursor.Decode(&m); err == nil {
			results = append(results, m)
		}
	}
	return results, nil
}

