package handlers

import (
	"contacts/db"
	"contacts/middlewares"
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User handler definition
type UsersHandler struct {
	Col db.CollectionAPI
}

// Handle users signup and validate request body
func (u *UsersHandler) Signup(c echo.Context) error {
	var user models.User
	c.Echo().Validator = &UsersValidator{validator: v}

	if err := c.Bind(&user); err != nil {
		return c.JSON(422, "Unable to parse request body")
	}

	if err := c.Validate(&user); err != nil {
		return c.JSON(400, "Invalid request body")
	}

	result, httpErr := db.CreateUser(context.Background(), user, u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, result)
}

// Handle users login and generate JWT token
func (u *UsersHandler) Login(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(422, "Unable to parse request body")
	}

	logUser, httpErr := db.LoginUser(context.Background(), user, u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	token, httpErr := logUser.GenerateToken()
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	c.Response().Header().Add("x-auth-token", "Bearer "+token)
	return c.JSON(200, models.User{Username: logUser.Username})
}

// Handle retrieve user info. Only return the email
func (u *UsersHandler) GetUser(c echo.Context) error {
	user, httpErr := db.RetrieveUser(context.Background(), c.Param("id"), u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, user)
}

// Get users followers
func (u *UsersHandler) GetFollowers(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(500, "Unable to convert to object id")
	}

	users, httpErr := db.GetUserFollowers(context.Background(), id, u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, users)
}

func (u *UsersHandler) GetUserPosts(c echo.Context) error {
	posts, httpErr := db.RetrievetUserPosts(context.Background(), c.Param("id"), u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, posts)
}

// Handle following users
func (u *UsersHandler) FollowUser(c echo.Context) error {
	toID := c.Param("id")
	fromID := userIDFromToken(c)

	if err := db.SetFollowUser(context.Background(), fromID, toID, u.Col); err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, "User followed successfuly")
}

// Get user id from token
func userIDFromToken(c echo.Context) string {
	_, claims := middlewares.GetToken(c)
	return claims["user_id"].(string)
}
