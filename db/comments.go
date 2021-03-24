package db

import (
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handle comments creation in the db
func CreateComment(ctx context.Context, id string, comment models.Comment, collection CollectionAPI) (models.Post, *echo.HTTPError) {
	var post models.Post
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, echo.NewHTTPError(400, "Unable to convert to object id")
	}
	filter := bson.M{"_id": docID}

	result := collection.FindOne(ctx, filter)
	if err = result.Decode(&post); err != nil {
		return post, echo.NewHTTPError(422, "Unable to parse retrieved post")
	}

	comment.ID = primitive.NewObjectID()
	post.Comments = append(post.Comments, comment)

	res, err := collection.UpdateOne(ctx, filter, bson.M{"$set": post})
	if err != nil {
		return post, echo.NewHTTPError(500, "Unable to find posts")
	}
	if res.ModifiedCount == 0 || res.MatchedCount == 0 {
		return post, echo.NewHTTPError(404, "Post does not exist")
	}

	return post, nil
}

// Handle Delete comments from the db
func RemoveComment(ctx context.Context, postID string, commentID string, collection CollectionAPI) (models.Post, *echo.HTTPError) {
	var post models.Post

	postDocID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return post, echo.NewHTTPError(400, "Unable to converto to object id")
	}

	commetDocID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return post, echo.NewHTTPError(400, "Unable to converto to object id")
	}

	filter := bson.M{"_id": postDocID}
	res := collection.FindOne(ctx, filter)
	if err = res.Decode(&post); err != nil {
		return post, echo.NewHTTPError(500, "Unable to decode retrieved post")
	}

	for i, comment := range post.Comments {
		if comment.ID == commetDocID {
			post.Comments = append(post.Comments[:i], post.Comments[i+1:]...)
			break
		}
	}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": post})
	if err != nil {
		return post, echo.NewHTTPError(500, "Unable to update the post")
	}

	return post, nil
}
