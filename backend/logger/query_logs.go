package logger

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// QueryLogs 查询日志，支持过滤和分页
func QueryLogs(ctx context.Context, filter bson.M, page, pageSize int) ([]bson.M, int64, error) {
	col := db.Collection("logs")
	findOpts := options.Find().SetSort(bson.M{"request_time": -1}).SetSkip(int64((page-1)*pageSize)).SetLimit(int64(pageSize))
	cursor, err := col.Find(ctx, filter, findOpts)
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
	total, _ := col.CountDocuments(ctx, filter)
	return results, total, nil
}

