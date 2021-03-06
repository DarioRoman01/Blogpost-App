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

// Handle users creation in the db
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

// Handle users login and check credentials
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

// check if the given password is correct
func isValidCredential(givenPWD, storePWD string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storePWD), []byte(givenPWD)); err != nil {
		return false
	}

	return true
}

// get user by id
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

// Manage the following system in the db fromID(requesting user) toID(users that requesting user want to follow)
// also check if the requesting user already follows userTo and perform follow or unfollow
func SetFollowUser(ctx context.Context, fromID string, toID string, collection CollectionAPI) *echo.HTTPError {
	var userFrom models.User
	var userTo models.User

	fromDocID, err := primitive.ObjectIDFromHex(fromID)
	if err != nil {
		return echo.NewHTTPError(500, "Unable to convert to object id")
	}

	toDocID, err := primitive.ObjectIDFromHex(toID)
	if err != nil {
		return echo.NewHTTPError(500, "Unable to convert to object id")
	}

	result := collection.FindOne(ctx, bson.M{"_id": fromDocID})
	if err = result.Decode(&userFrom); err != nil {
		return echo.NewHTTPError(500, "Unable to decode retrieved user")
	}

	result = collection.FindOne(ctx, bson.M{"_id": toDocID})
	if err = result.Decode(&userTo); err != nil {
		return echo.NewHTTPError(500, "Unable to decode retrieved user")
	}

	if !contains(userTo.Followers, fromID) {
		userFrom.Following = append(userFrom.Following, toID)
		userTo.Followers = append(userTo.Followers, fromID)

		_, err := collection.UpdateOne(ctx, bson.M{"_id": fromDocID}, bson.M{"$set": userFrom})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update user info")
		}

		_, err = collection.UpdateOne(ctx, bson.M{"_id": toDocID}, bson.M{"$set": userTo})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update user data")
		}
	} else {
		_, err = collection.UpdateOne(ctx, bson.M{"_id": toDocID}, bson.M{"$pull": bson.M{"followers": fromDocID}})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update user")
		}

		_, err = collection.UpdateOne(ctx, bson.M{"_id": fromDocID}, bson.M{"$pull:": bson.M{"following": toDocID}})
		if err != nil {
			return echo.NewHTTPError(500, "Unable to update user")
		}
	}

	return nil
}

// Retrieve all followers of the user
func GetUserFollowers(ctx context.Context, id primitive.ObjectID, collection CollectionAPI) ([]models.User, *echo.HTTPError) {
	var followers []primitive.ObjectID
	var users []models.User
	var user models.User

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Decode(&user); err != nil {
		return nil, echo.NewHTTPError(500, "Unable to decode retrieved user")
	}

	for _, userId := range user.Followers {
		docID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return nil, echo.NewHTTPError(500, "Unable to convert to object id")
		}
		followers = append(followers, docID)
	}

	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": followers}})
	if err != nil {
		return nil, echo.NewHTTPError(404, "Unable to find users")
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, echo.NewHTTPError(500, "Unable to parse retrieved users")
	}

	return users, nil
}
