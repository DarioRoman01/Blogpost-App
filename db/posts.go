package db

import (
	"contacts/models"
	"context"
	"encoding/json"
	"io"

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
		return post, echo.NewHTTPError(400, "Unable to convert to object id")
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
		return posts, echo.NewHTTPError(404, "Unable to find posts")
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
		return nil, echo.NewHTTPError(400, "Unable to convert to object id")
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		return nil, echo.NewHTTPError(500, "unable to delete post")
	}

	if result.DeletedCount == 0 {
		return nil, echo.NewHTTPError(404, "Post id does not exist")
	}

	return result, nil
}

// Handle update post data
func UpdatePost(ctx context.Context, id string, reqBody io.ReadCloser, collection CollectionAPI) (models.Post, *echo.HTTPError) {
	var post models.Post

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, echo.NewHTTPError(400, "Unable to convert to object id")
	}
	filter := bson.M{"_id": docID}
	result := collection.FindOne(ctx, filter)

	if err = result.Decode(&post); err != nil {
		return post, echo.NewHTTPError(404, "Post not found")
	}

	if err := json.NewDecoder(reqBody).Decode(&post); err != nil {
		return post, echo.NewHTTPError(422, "Unable to parse request payload")
	}

	if _, err = collection.UpdateOne(ctx, filter, bson.M{"$set": post}); err != nil {
		return post, echo.NewHTTPError(500, "Unable to update post")
	}

	return post, nil
}

// Retrieve all posts from one user
func RetrievetUserPosts(ctx context.Context, id string, collection CollectionAPI) ([]models.Post, *echo.HTTPError) {
	var posts []models.Post

	cursor, err := collection.Find(ctx, bson.M{"from": id})
	if err != nil {
		return nil, echo.NewHTTPError(404, "Unable to find posts")
	}

	if err = cursor.All(ctx, &posts); err != nil {
		return nil, echo.NewHTTPError(500, "Unable to decode retrieved posts")
	}

	return posts, nil
}

// check if requesting user already like the post and remove or add the like to the post
func SetLike(ctx context.Context, userID, postID string, collection CollectionAPI) *echo.HTTPError {
	var post models.Post

	docID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return echo.NewHTTPError(500, "Unable to convert to object id")
	}

	result := collection.FindOne(ctx, bson.M{"_id": docID})
	if result.Err() != nil {
		return echo.NewHTTPError(400, "Post does not exist")
	}

	if err = result.Decode(&post); err != nil {
		return echo.NewHTTPError(500, "uanble to decode retrieved post")
	}

	if !contains(post.LikedBy, userID) {
		post.LikedBy = append(post.LikedBy, userID)
		post.Likes++

		_, err = collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{"$set": post})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update post")
		}
	} else {
		index := GetIndex(post.LikedBy, userID)
		post.LikedBy = append(post.LikedBy[:index], post.LikedBy[index+1:]...)
		post.Likes--
		_, err = collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.M{"$set": post})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update post")
		}
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetIndex(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return 0
}
