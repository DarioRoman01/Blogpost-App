package handlers

import (
	"contacts/db"
	"contacts/middlewares"
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
)

type UsersHandler struct {
	Col db.CollectionAPI
}

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

func (u *UsersHandler) GetUser(c echo.Context) error {
	user, httpErr := db.RetrieveUser(context.Background(), c.Param("id"), u.Col)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, user)
}

func (u *UsersHandler) FollowUser(c echo.Context) error {
	toID := c.Param("id")
	fromID := userIDFromToken(c)

	if err := db.SetFollowUser(context.Background(), fromID, toID, u.Col); err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, "User followed successfuly")
}

func userIDFromToken(c echo.Context) string {
	_, claims := middlewares.GetToken(c)
	return claims["user_id"].(string)
}
