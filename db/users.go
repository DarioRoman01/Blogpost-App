package db

import (
	"contacts/models"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx context.Context, user models.User, collection CollectionAPI) (*mongo.InsertOneResult, *echo.HTTPError) {
	var newUser models.User

	result := collection.FindOne(ctx, bson.M{"username": user.Username, "email": user.Email})
	err := result.Decode(&newUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, echo.NewHTTPError(500, "Unable to decode new user")
	}

	if newUser.Email != "" || newUser.Username != "" {
		return nil, echo.NewHTTPError(400, "That user already exist")
	}

	hashpwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, echo.NewHTTPError(500, "Unable to hash password")
	}

	user.ID = primitive.NewObjectID()
	user.Password = string(hashpwd)

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, echo.NewHTTPError(400, "That user already exist")
	}

	return res, nil
}

func LoginUser(ctx context.Context, reqUser models.User, collection CollectionAPI) (models.User, *echo.HTTPError) {
	var user models.User

	result := collection.FindOne(ctx, bson.M{"username": reqUser.Username})
	err := result.Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return reqUser, echo.NewHTTPError(400, "Unable to parse request user")
	}

	if err == mongo.ErrNoDocuments {
		return reqUser, echo.NewHTTPError(400, "Users do not exist")
	}

	if !isValidCredential(reqUser.Password, user.Password) {
		return reqUser, echo.NewHTTPError(400, "Invalid credentials")
	}

	return user, nil
}

func isValidCredential(givenPWD, storePWD string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storePWD), []byte(givenPWD)); err != nil {
		return false
	}

	return true
}

func RetrieveUser(ctx context.Context, id string, collection CollectionAPI) (models.User, *echo.HTTPError) {
	var user models.User

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, echo.NewHTTPError(500, "Unable to convert to object id")
	}

	result := collection.FindOne(ctx, bson.M{"_id": docID})
	if err = result.Decode(&user); err != nil {
		return user, echo.NewHTTPError(422, "Unable to parse retrieved user")
	}

	return models.User{Username: user.Username}, nil
}

func SetFollowUser(ctx context.Context, fromID string, toID string, collection CollectionAPI) *echo.HTTPError {
	fromDocID, err := primitive.ObjectIDFromHex(fromID)
	if err != nil {
		return echo.NewHTTPError(500, "Unable to convert to object id")
	}

	toDocID, err := primitive.ObjectIDFromHex(toID)
	if err != nil {
		return echo.NewHTTPError(500, "Unable to convert to object id")
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": toDocID}, bson.M{"$addToSet": bson.M{"followers": fromDocID}})
	if err != nil {
		return echo.NewHTTPError(500, "Unable to update user")
	}

	if res.MatchedCount == 0 {
		return echo.NewHTTPError(400, "User id does not exist")
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": fromDocID}, bson.M{"$addToSet": bson.M{"following": toDocID}})
	if err != nil {
		return echo.NewHTTPError(500, "Unable to update user")
	}

	return nil
}
