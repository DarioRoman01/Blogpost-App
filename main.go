package main

import (
	"contacts/config"
	"contacts/db"
	"contacts/handlers"
	"contacts/middlewares"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	usersColl *mongo.Collection
	postsColl *mongo.Collection
	cfg       config.Properties
)

// init get connection with the db and read config
func init() {
	usersColl, postsColl = db.GetConnection()
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Unable to read configuration")
	}
}

func main() {
	// create new echo instance and set middlewares
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middlewares.LoggerMiddleware())
	e.Use(middlewares.JwtMiddleware())

	// instance handlers uh(users handler) ph(posts handlers)
	uh := &handlers.UsersHandler{Col: usersColl}
	ph := &handlers.PostsHandler{Col: postsColl}

	// posts endpoints
	e.POST("/posts/create", ph.CreatePost)
	e.GET("/posts/:id", ph.GetPost)
	e.GET("/posts", ph.ListPosts)
	e.DELETE("/posts/:id", ph.RemovePost, middlewares.IsPostOwner)
	e.PATCH("/posts/:id", ph.PostUpdate, middlewares.IsPostOwner)

	// users endpoints
	e.POST("/users/signup", uh.Signup)
	e.POST("/users/login", uh.Login)
	e.GET("/users/:id", uh.GetUser)
	e.POST("/users/:id/follow", uh.FollowUser)

	// initializer server
	e.Logger.Info("Listening on port %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
