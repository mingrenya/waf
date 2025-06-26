package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"coraza-waf/backend/models"
)

type RuleRepository interface {
	CreateRule(ctx context.Context, rule *models.Rule) error
	UpdateRule(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteRule(ctx context.Context, id primitive.ObjectID) error
	GetRule(ctx context.Context, id primitive.ObjectID) (*models.Rule, error)
	ListRules(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Rule, int64, error)
	EnableRule(ctx context.Context, id primitive.ObjectID, enabled bool) error
}

type MongoRuleRepo struct {
	col *mongo.Collection
}

func NewRuleRepository(col *mongo.Collection) RuleRepository {
	return &MongoRuleRepo{col: col}
}

func (r *MongoRuleRepo) CreateRule(ctx context.Context, rule *models.Rule) error {
	rule.ID = primitive.NewObjectID()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, rule)
	return err
}

func (r *MongoRuleRepo) UpdateRule(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.col.UpdateByID(ctx, id, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("rule not found")
	}
	return nil
}

func (r *MongoRuleRepo) DeleteRule(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("rule not found")
	}
	return nil
}

func (r *MongoRuleRepo) GetRule(ctx context.Context, id primitive.ObjectID) (*models.Rule, error) {
	var rule models.Rule
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&rule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("rule not found")
		}
		return nil, err
	}
	return &rule, nil
}

func (r *MongoRuleRepo) ListRules(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Rule, int64, error) {
	findOpts := options.Find().SetSort(bson.M{"priority": 1, "created_at": -1}).SetSkip(int64((page-1)*pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var rules []models.Rule
	for cursor.Next(ctx) {
		var rule models.Rule
		if err := cursor.Decode(&rule); err == nil {
			rules = append(rules, rule)
		}
	}
	total, _ := r.col.CountDocuments(ctx, filter)
	return rules, total, nil
}

func (r *MongoRuleRepo) EnableRule(ctx context.Context, id primitive.ObjectID, enabled bool) error {
	return r.UpdateRule(ctx, id, bson.M{"enabled": enabled})
}
