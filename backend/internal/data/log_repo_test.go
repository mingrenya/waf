package data

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getTestMongoCollection(t *testing.T) *mongo.Collection {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Disconnect(context.Background())
	})
	db := client.Database("waf_test")
	col := db.Collection("logs_test")
	_ = col.Drop(context.Background())
	return col
}

func TestMongoLogRepo_Integration(t *testing.T) {
	col := getTestMongoCollection(t)
	repo := NewLogRepository(col)
	ctx := context.Background()

	logDoc := bson.M{"foo": "bar", "request_time": time.Now()}

	t.Run("InsertLog", func(t *testing.T) {
		err := repo.InsertLog(ctx, logDoc)
		if err != nil {
			t.Errorf("InsertLog failed: %v", err)
		}
	})

	t.Run("QueryLogs", func(t *testing.T) {
		logs, total, err := repo.QueryLogs(ctx, bson.M{"foo": "bar"}, 1, 10)
		if err != nil {
			t.Errorf("QueryLogs failed: %v", err)
		}
		if len(logs) == 0 || total < 1 {
			t.Errorf("QueryLogs returned no results or invalid total")
		}
	})

	t.Run("Aggregate", func(t *testing.T) {
		pipeline := mongo.Pipeline{
			bson.D{{"$match", bson.D{{"foo", "bar"}}}},
			bson.D{{"$group", bson.D{{"_id", "$foo"}, {"count", bson.D{{"$sum", 1}}}}}},
		}
		res, err := repo.Aggregate(ctx, pipeline)
		if err != nil {
			t.Errorf("Aggregate failed: %v", err)
		}
		if len(res) == 0 {
			t.Errorf("Aggregate returned no results")
		}
	})

	t.Run("DeleteLogs", func(t *testing.T) {
		count, err := repo.DeleteLogs(ctx, bson.M{"foo": "bar"})
		if err != nil {
			t.Errorf("DeleteLogs failed: %v", err)
		}
		if count < 1 {
			t.Errorf("DeleteLogs returned invalid count")
		}
	})
}

