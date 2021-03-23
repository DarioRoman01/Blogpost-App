package handlers

import "gopkg.in/go-playground/validator.v9"

var v = validator.New()

type UsersValidator struct {
	validator *validator.Validate
}

// Validate users definition
func (u *UsersValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}

type PostsValidator struct {
	validator *validator.Validate
}

// Validate Posts definition
func (p *PostsValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

type CommentValidator struct {
	validator *validator.Validate
}

func (c *CommentValidator) Validate(i interface{}) error {
	return c.validator.Struct(i)
}
