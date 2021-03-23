package models

import (
	"contacts/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var cfg config.Properties

// User definition
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required,min=3"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required,min=8,max=300"`
	Followers []string           `json:"Followers,omitempty" bson:"followers,omitempty"`
	Following []string           `json:"following,omitempty" bson:"following,omitempty"`
}

// util to function to generate token for requesting user
func (u User) GenerateToken() (string, *echo.HTTPError) {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Cannot read configuration")
	}

	claims := jwt.MapClaims{}
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(cfg.JwtTokenSecret))
	if err != nil {
		return "", echo.NewHTTPError(500, "Unable to create token")
	}

	return token, nil
}
