package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog 操作审计日志：记录写操作的用户、动作、实体与请求上下文。
type AuditLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username   string             `bson:"username" json:"username"`
	Action     string             `bson:"action" json:"action"`
	Entity     string             `bson:"entity" json:"entity"`
	EntityID   string             `bson:"entity_id" json:"entity_id"`
	Method     string             `bson:"method" json:"method"`
	Path       string             `bson:"path" json:"path"`
	StatusCode int                `bson:"status_code" json:"status_code"`
	IP         string             `bson:"ip" json:"ip"`
	RequestID  string             `bson:"request_id" json:"request_id"`
	Detail     string             `bson:"detail" json:"detail"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
