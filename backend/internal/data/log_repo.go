package data

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
	findOpts := options.Find().SetSort(bson.M{"request_time": -1}).SetSkip(int64((page-1)*pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
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

func (r *MongoLogRepo) InsertLog(ctx context.Context, log interface{}) error {
	_, err := r.col.InsertOne(ctx, log)
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
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
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

