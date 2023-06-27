package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// LoginRequestBody defines the incoming login JSON body. Implements chi render.Binder
type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (schema *LoginRequestBody) Bind(r *http.Request) error {
	return validator.New().Struct(schema)
}

// SignupJSONBody defines the incoming signup JSON body. Implements chi render.Binder
type SignupJSONBody struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
	Password  string `json:"password"  validate:"required,min=8"`
	Email     string `json:"email"  validate:"required,email"`
}

func (schema *SignupJSONBody) Bind(r *http.Request) error {
	return validator.New().Struct(schema)
}
