package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserId int64    `bson:"user_id"`
	URLs   []string `bson:"urls"`
}

type UserResponse struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserId int64              `bson:"user_id"`
	URLs   []string           `bson:"urls"`
}
