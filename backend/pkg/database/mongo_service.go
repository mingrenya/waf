package database

import (
	"context"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"coraza-waf/backend/pkg/logging"
)

type MongoService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoService(uri, dbName, collName string) (*MongoService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	coll := client.Database(dbName).Collection(collName)
	return &MongoService{client: client, collection: coll}, nil
}

func (m *MongoService) InsertLog(logEntry *logging.WafLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Mongo InsertOne error: %v", err)
		return err
	}
	return nil
}

