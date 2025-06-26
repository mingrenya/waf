package database

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoService 封装 MongoDB 操作
type MongoService struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewMongoService 初始化 MongoService，参数：uri, 数据库名, 集合名
func NewMongoService(uri, dbName, collectionName string) (*MongoService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		cancel()
		return nil, err
	}

	// Ping 确认连接
	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoService{
		client:     client,
		collection: collection,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// InsertLog 插入日志文档
func (m *MongoService) InsertLog(document interface{}) error {
	_, err := m.collection.InsertOne(m.ctx, document)
	return err
}

// Close 断开连接，释放资源
func (m *MongoService) Close() error {
	m.cancel()
	return m.client.Disconnect(m.ctx)
}

