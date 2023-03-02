package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	database = "pingrobot"
	url      = "mongodb+srv://qara-qurt:Serikov12@pingrobot.fyvyccb.mongodb.net/?retryWrites=true&w=majority"
)

type Storage struct {
	DB    *mongo.Database
	users *mongo.Collection
}

func InitDatabase() (*Storage, error) {
	// Create a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// Connect to the MongoDB database
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Access the MongoDB database
	db := client.Database(database)
	fmt.Println("Connected to MongoDB!")
	return &Storage{
		DB:    db,
		users: db.Collection("users"),
	}, err
}

func (s Storage) Create(userId int64) error {
	user := User{UserId: userId, URLs: []string{}}
	collection := s.DB.Collection("users")
	if _, err := collection.InsertOne(context.Background(), user); err != nil {
		return err
	}
	return nil
}

func (s Storage) Add(url string, userId int64) error {
	user, err := s.Get(userId)
	if err != nil {
		return err
	}

	newURLs := append(user.URLs, url)
	user.URLs = newURLs
	update := bson.D{{"$set", user}}

	filter := bson.D{{"user_id", userId}}
	if _, err := s.users.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}

func (s Storage) Get(userId int64) (UserResponse, error) {
	filter := bson.D{{"user_id", userId}}

	var res UserResponse
	user := s.users.FindOne(context.Background(), filter)
	if user.Err() != nil {
		return UserResponse{}, user.Err()
	}
	user.Decode(&res)
	return res, nil
}
