package middlewares

import (
	"contacts/config"
	"contacts/db"
	"contacts/models"
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var cfg config.Properties

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Unable to read configuration")
	}

}

func LoggerMiddleware() echo.MiddlewareFunc {
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${host} ${method} ${status} ${uri} ${user_agent}` +
			`${status} ${error} ${latency_human}` + "\n",
	})

	return logger
}

func JwtMiddleware() echo.MiddlewareFunc {
	jwtMidd := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.JwtTokenSecret),
		TokenLookup: "header:x-auth-token",
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/users/login" || c.Path() == "/users/signup" {
				return true
			}
			return false
		},
	})

	return jwtMidd
}

func IsPostOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var post models.Post

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			echo.NewHTTPError(500, "Unable to convert to object id")
		}

		_, postColl := db.GetConnection()
		defer postColl.Database().Client().Disconnect(context.Background())

		result := postColl.FindOne(context.Background(), bson.M{"_id": id})
		if err := result.Decode(&post); err != nil {
			return echo.NewHTTPError(500, "Unable to decode retrieved user")
		}

		_, claims := GetToken(c)

		if claims["user_id"] != post.From {
			return echo.NewHTTPError(403, "You dont have permissions to perform this action")
		}

		return next(c)
	}

}

func GetToken(c echo.Context) (*jwt.Token, jwt.MapClaims) {
	headerToken := c.Request().Header.Get("x-auth-token")
	strToken := strings.Split(headerToken, " ")[1]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtTokenSecret), nil
	})
	if err != nil {
		return nil, nil
	}

	return token, claims
}
