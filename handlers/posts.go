package handlers

import (
	"contacts/db"
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post handler definition
type PostsHandler struct {
	Col db.CollectionAPI
}

// Handle requesting data and validation for posts creation
func (p *PostsHandler) CreatePost(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(userIDFromToken(c))
	if err != nil {
		return c.JSON(500, "Unable to convert to object")
	}

	var post models.Post
	post.From = id.Hex()

	c.Echo().Validator = &PostsValidator{validator: v}

	if err := c.Bind(&post); err != nil {
		return c.JSON(422, "Unable to parse request body")
	}

	if err := c.Validate(&post); err != nil {
		return c.JSON(400, "Invalid request body")
	}

	result, httpErr := db.InsertPost(context.Background(), post, p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, result)
}

// retrieve one post
func (p *PostsHandler) GetPost(c echo.Context) error {
	post, httpErr := db.FindPost(context.Background(), c.Param("id"), p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, post)
}

// list posts based on users that requesting user is following
func (p *PostsHandler) ListPosts(c echo.Context) error {
	var user models.User
	id := userIDFromToken(c)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(500, "Unable to converto to object id")
	}

	ctx := context.Background()

	usersColl, _ := db.GetConnection()
	result := usersColl.FindOne(ctx, bson.M{"_id": docID})
	if err := result.Decode(&user); err != nil {
		return c.JSON(500, "Something wrong happend in the request")
	}

	defer usersColl.Database().Client().Disconnect(ctx)
	user.Following = append(user.Following, id)
	res, httpErr := db.FindPosts(ctx, user.Following, p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, res)
}

// handle delete product request
func (p *PostsHandler) RemovePost(c echo.Context) error {
	delIDS, httpErr := db.DeletePost(context.Background(), c.Param("id"), p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, delIDS)
}

// Handle post update request
func (p *PostsHandler) PostUpdate(c echo.Context) error {
	post, httpErr := db.UpdatePost(context.Background(), c.Param("id"), c.Request().Body, p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, post)
}

// Handle validation and parsing request data for comments creation
func (p *PostsHandler) CommentPost(c echo.Context) error {
	id := userIDFromToken(c)
	c.Echo().Validator = &CommentValidator{validator: v}
	var comment models.Comment
	comment.From = id

	if err := c.Bind(&comment); err != nil {
		return c.JSON(422, "Unable to parse request body")
	}

	if err := c.Validate(&comment); err != nil {
		return c.JSON(400, "Invalid request body")
	}

	result, httpErr := db.CreateComment(context.Background(), c.Param("id"), comment, p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, result)
}

// Handle delete comment request
func (p *PostsHandler) DeleteComment(c echo.Context) error {
	postID := c.Param("id")
	commentID := c.Param("cid")

	post, err := db.RemoveComment(context.Background(), postID, commentID, p.Col)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, post)
}

func (p *PostsHandler) ToggleLikePost(c echo.Context) error {
	postID := c.Param("id")
	userID := userIDFromToken(c)

	httpErr := db.SetLike(context.Background(), userID, postID, p.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, "request was successfully")
}
