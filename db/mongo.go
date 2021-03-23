package db

import (
	"contacts/config"
	"context"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection interface
type CollectionAPI interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

func GetConnection() (*mongo.Collection, *mongo.Collection) {
	var cfg config.Properties
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Unable to read configuration")
	}

	ctx := context.Background()

	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
	if err != nil {
		panic("Unable to connect to mongo")
	}

	db := client.Database(cfg.DBName)
	usersCollection := db.Collection(cfg.UsersCollection)
	postsCollection := db.Collection(cfg.PostsCollection)

	isUsernameUnique := true
	usernameIndexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: &options.IndexOptions{
			Unique: &isUsernameUnique,
		},
	}

	isEmailUnique := true
	emailIndexModel := mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: &options.IndexOptions{
			Unique: &isEmailUnique,
		},
	}

	_, err = db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{usernameIndexModel, emailIndexModel})
	if err != nil {
		panic("Unable to create indexes")
	}

	return usersCollection, postsCollection
}
