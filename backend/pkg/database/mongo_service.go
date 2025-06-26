package database

import (
	"context"
	//"fmt"
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
	// 创建上下文用于连接和后续操作
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // [6] 延迟关闭上下文以避免资源泄漏

	// 配置 MongoDB 连接参数
	clientOptions := options.Client().ApplyURI(uri) // [4] 生成 ClientOptions 实例
	client, err := mongo.Connect(clientOptions)    // [1] 仅传递 options.ClientOptions 参数
	if err != nil {
			cancel()
		return nil, err
	}

	// 使用 context.TODO() 作为 Ping 的上下文 [6]
	if err := client.Ping(ctx, nil); err != nil {
			cancel() // [2] 失败时取消上下文
			return nil, err
	}

	// 获取数据库和集合
	collection := client.Database(dbName).Collection(collectionName)

	return &MongoService{
		client:     client,
		collection: collection,
		ctx:        ctx, // [6] 保留上下文用于后续操作（如 InsertOne、Find 等）
		cancel:     cancel, // [2] 保留 cancel 用于主动关闭连接
	}, nil
}

// InsertLog 插入日志文档
func (m *MongoService) InsertLog(document interface{}) error {
	_, err := m.collection.InsertOne(m.ctx, document) // [6] 正确传递 context.Context
	return err
}

// Close 断开连接，释放资源
func (m *MongoService) Close() error {
	m.cancel() // [2] 主动取消上下文以触发清理操作
	if err := m.client.Disconnect(m.ctx); err != nil { // [2] 传递 context.Context 作为 Disconnect 参数
		return err
	}
	return nil
}

