package db

import (
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// insert the post in the db
func InsertPost(ctx context.Context, post models.Post, collection CollectionAPI) (*mongo.InsertOneResult, *echo.HTTPError) {
	post.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, echo.NewHTTPError(422, "Unable to cretae post")
	}

	return result, nil
}

// retrieve one post
func FindPost(ctx context.Context, id string, collection CollectionAPI) (models.Post, *echo.HTTPError) {
	var post models.Post

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, echo.NewHTTPError(500, "Unable to convert to object id")
	}

	result := collection.FindOne(ctx, bson.M{"_id": docID})
	if err = result.Decode(&post); err != nil {
		return post, echo.NewHTTPError(404, "Post not found")
	}

	return post, nil
}

// Get posts based on the users that the requesting user is following
func FindPosts(ctx context.Context, follows []string, collection CollectionAPI) ([]models.Post, *echo.HTTPError) {
	var posts []models.Post

	cursor, err := collection.Find(ctx, bson.M{"from": bson.M{"$in": follows}})
	if err != nil {
		return posts, echo.NewHTTPError(500, "Unable to find posts")
	}

	if err = cursor.All(ctx, &posts); err != nil {
		return posts, echo.NewHTTPError(500, "Unable to parse retrieved posts")
	}

	return posts, nil
}

// Delete post from the db
func DeletePost(ctx context.Context, id string, collection CollectionAPI) (*mongo.DeleteResult, *echo.HTTPError) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Unable to convert to object id")
	}

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": docID})
	if err != nil {
		return nil, echo.NewHTTPError(500, "unable to delete post")
	}

	if result.DeletedCount == 0 {
		return nil, echo.NewHTTPError(400, "Post id does not exist")
	}

	return result, nil
}
