package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseObjectID 将字符串解析为 primitive.ObjectID。
func ParseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
