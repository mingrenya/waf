package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Rule struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Content      string             `bson:"content" json:"content"`
	RuleType     string             `bson:"rule_type" json:"rule_type"`
	Enabled      bool               `bson:"enabled" json:"enabled"`
	Tags         []string           `bson:"tags" json:"tags"`
	Priority     int                `bson:"priority" json:"priority"`
	Description  string             `bson:"description" json:"description"`
	Version      string             `bson:"version" json:"version"`
	Source       string             `bson:"source" json:"source"`
	Action       string             `bson:"action" json:"action"`
	Status       string             `bson:"status" json:"status"`
	UpdateBy     string             `bson:"update_by" json:"update_by"`
	UpdateReason string             `bson:"update_reason" json:"update_reason"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	Format       string             `bson:"format" json:"format"` // 规则内容格式：modsec/json/expr等
}
